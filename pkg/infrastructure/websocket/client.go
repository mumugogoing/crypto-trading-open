package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// MessageHandler handles incoming WebSocket messages
type MessageHandler func(message []byte) error

// ErrorHandler handles WebSocket errors
type ErrorHandler func(err error)

// Client represents a WebSocket client
type Client struct {
	url             string
	conn            *websocket.Conn
	logger          *zap.Logger
	messageHandler  MessageHandler
	errorHandler    ErrorHandler
	reconnect       bool
	reconnectDelay  time.Duration
	pingInterval    time.Duration
	readTimeout     time.Duration
	writeTimeout    time.Duration
	mu              sync.RWMutex
	connected       bool
	stopChan        chan struct{}
	reconnectChan   chan struct{}
	writeChan       chan []byte
}

// Config represents WebSocket client configuration
type Config struct {
	URL            string
	ReconnectDelay time.Duration
	PingInterval   time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	EnableReconnect bool
}

// NewClient creates a new WebSocket client
func NewClient(config *Config, logger *zap.Logger) *Client {
	return &Client{
		url:            config.URL,
		logger:         logger,
		reconnect:      config.EnableReconnect,
		reconnectDelay: config.ReconnectDelay,
		pingInterval:   config.PingInterval,
		readTimeout:    config.ReadTimeout,
		writeTimeout:   config.WriteTimeout,
		stopChan:       make(chan struct{}),
		reconnectChan:  make(chan struct{}, 1),
		writeChan:      make(chan []byte, 100),
	}
}

// Connect connects to the WebSocket server
func (c *Client) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connected {
		return fmt.Errorf("already connected")
	}

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.DialContext(ctx, c.url, nil)
	if err != nil {
		c.logger.Error("Failed to connect to WebSocket", zap.String("url", c.url), zap.Error(err))
		return err
	}

	c.conn = conn
	c.connected = true
	c.logger.Info("Connected to WebSocket", zap.String("url", c.url))

	// Start read/write loops
	go c.readLoop()
	go c.writeLoop()
	go c.pingLoop()

	return nil
}

// Disconnect disconnects from the WebSocket server
func (c *Client) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.connected {
		return nil
	}

	close(c.stopChan)
	c.connected = false

	if c.conn != nil {
		err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			c.logger.Warn("Failed to send close message", zap.Error(err))
		}
		c.conn.Close()
		c.conn = nil
	}

	c.logger.Info("Disconnected from WebSocket")
	return nil
}

// IsConnected returns whether the client is connected
func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}

// SendMessage sends a message to the WebSocket server
func (c *Client) SendMessage(message []byte) error {
	if !c.IsConnected() {
		return fmt.Errorf("not connected")
	}

	select {
	case c.writeChan <- message:
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("send timeout")
	}
}

// SendJSON sends a JSON message to the WebSocket server
func (c *Client) SendJSON(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return c.SendMessage(data)
}

// SetMessageHandler sets the message handler
func (c *Client) SetMessageHandler(handler MessageHandler) {
	c.messageHandler = handler
}

// SetErrorHandler sets the error handler
func (c *Client) SetErrorHandler(handler ErrorHandler) {
	c.errorHandler = handler
}

// readLoop reads messages from the WebSocket
func (c *Client) readLoop() {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Error("Read loop panic", zap.Any("panic", r))
		}
		c.handleDisconnect()
	}()

	for {
		select {
		case <-c.stopChan:
			return
		default:
		}

		if c.readTimeout > 0 {
			c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
		}

		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Error("WebSocket read error", zap.Error(err))
			}
			if c.errorHandler != nil {
				c.errorHandler(err)
			}
			return
		}

		if c.messageHandler != nil {
			if err := c.messageHandler(message); err != nil {
				c.logger.Error("Message handler error", zap.Error(err))
			}
		}
	}
}

// writeLoop writes messages to the WebSocket
func (c *Client) writeLoop() {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Error("Write loop panic", zap.Any("panic", r))
		}
	}()

	for {
		select {
		case <-c.stopChan:
			return
		case message := <-c.writeChan:
			if c.writeTimeout > 0 {
				c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				c.logger.Error("WebSocket write error", zap.Error(err))
				if c.errorHandler != nil {
					c.errorHandler(err)
				}
				return
			}
		}
	}
}

// pingLoop sends periodic ping messages
func (c *Client) pingLoop() {
	if c.pingInterval == 0 {
		return
	}

	ticker := time.NewTicker(c.pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopChan:
			return
		case <-ticker.C:
			if !c.IsConnected() {
				return
			}
			if c.writeTimeout > 0 {
				c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.logger.Error("Failed to send ping", zap.Error(err))
				return
			}
		}
	}
}

// handleDisconnect handles disconnection and reconnection
func (c *Client) handleDisconnect() {
	c.mu.Lock()
	c.connected = false
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
	c.mu.Unlock()

	if c.reconnect {
		select {
		case c.reconnectChan <- struct{}{}:
			go c.reconnectLoop()
		default:
		}
	}
}

// reconnectLoop attempts to reconnect
func (c *Client) reconnectLoop() {
	for {
		select {
		case <-c.stopChan:
			return
		case <-c.reconnectChan:
			c.logger.Info("Attempting to reconnect...")
			time.Sleep(c.reconnectDelay)

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			err := c.Connect(ctx)
			cancel()

			if err != nil {
				c.logger.Error("Reconnect failed", zap.Error(err))
				select {
				case c.reconnectChan <- struct{}{}:
				default:
				}
			} else {
				c.logger.Info("Reconnected successfully")
				return
			}
		}
	}
}
