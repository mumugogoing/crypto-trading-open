package exchanges

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

// ExchangeStatus represents exchange status
type ExchangeStatus string

const (
	ExchangeStatusDisconnected  ExchangeStatus = "disconnected"
	ExchangeStatusConnecting    ExchangeStatus = "connecting"
	ExchangeStatusConnected     ExchangeStatus = "connected"
	ExchangeStatusAuthenticated ExchangeStatus = "authenticated"
	ExchangeStatusError         ExchangeStatus = "error"
	ExchangeStatusMaintenance   ExchangeStatus = "maintenance"
)

// ExchangeConfig represents exchange configuration
type ExchangeConfig struct {
	// Basic configuration
	ExchangeID   string
	Name         string
	ExchangeType ExchangeType

	// Authentication
	APIKey         string
	APISecret      string
	APIPassphrase  string
	WalletAddress  string

	// Network configuration
	Testnet bool
	BaseURL string
	WsURL   string

	// Trading configuration
	DefaultLeverage   int
	DefaultMarginMode string
	SymbolMapping     map[string]string

	// Limits configuration
	RateLimits map[string]interface{}
	Precision  map[string]interface{}

	// Feature flags
	EnableWebsocket      bool
	EnableAutoReconnect  bool
	EnableHeartbeat      bool

	// Timeout configuration
	ConnectTimeout    int
	RequestTimeout    int
	HeartbeatInterval int

	// Retry configuration
	MaxRetryAttempts int
	RetryDelay       float64

	// Extra parameters
	ExtraParams map[string]interface{}
}

// NewExchangeConfig creates a new ExchangeConfig with defaults
func NewExchangeConfig(exchangeID, name string, exchangeType ExchangeType, apiKey, apiSecret string) *ExchangeConfig {
	return &ExchangeConfig{
		ExchangeID:           exchangeID,
		Name:                 name,
		ExchangeType:         exchangeType,
		APIKey:               apiKey,
		APISecret:            apiSecret,
		Testnet:              false,
		DefaultLeverage:      1,
		DefaultMarginMode:    "cross",
		SymbolMapping:        make(map[string]string),
		RateLimits:           make(map[string]interface{}),
		Precision:            make(map[string]interface{}),
		EnableWebsocket:      true,
		EnableAutoReconnect:  true,
		EnableHeartbeat:      true,
		ConnectTimeout:       30,
		RequestTimeout:       10,
		HeartbeatInterval:    30,
		MaxRetryAttempts:     3,
		RetryDelay:           1.0,
		ExtraParams:          make(map[string]interface{}),
	}
}

// EventCallback is a function that handles exchange events
type EventCallback func(event interface{})

// ExchangeInterface defines the unified interface for all exchanges
type ExchangeInterface interface {
	// Lifecycle management
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	Authenticate(ctx context.Context) error
	HealthCheck(ctx context.Context) (map[string]interface{}, error)

	// Market data interface
	GetExchangeInfo(ctx context.Context) (*ExchangeInfo, error)
	GetTicker(ctx context.Context, symbol string) (*TickerData, error)
	GetOrderBook(ctx context.Context, symbol string, depth int) (*OrderBookData, error)
	GetTrades(ctx context.Context, symbol string, limit int) ([]*TradeData, error)
	GetOHLCV(ctx context.Context, symbol, timeframe string, since time.Time, limit int) ([]*OHLCVData, error)

	// Trading interface
	CreateOrder(ctx context.Context, symbol string, side OrderSide, orderType OrderType, amount, price decimal.Decimal, params map[string]interface{}) (*OrderData, error)
	CancelOrder(ctx context.Context, orderID, symbol string) (*OrderData, error)
	CancelAllOrders(ctx context.Context, symbol string) ([]*OrderData, error)
	GetOrder(ctx context.Context, orderID, symbol string) (*OrderData, error)
	GetOpenOrders(ctx context.Context, symbol string) ([]*OrderData, error)
	GetClosedOrders(ctx context.Context, symbol string, since time.Time, limit int) ([]*OrderData, error)

	// Account interface
	GetBalance(ctx context.Context) (map[string]*BalanceData, error)
	GetPositions(ctx context.Context, symbol string) ([]*PositionData, error)
	GetPosition(ctx context.Context, symbol string) (*PositionData, error)

	// Position management
	SetLeverage(ctx context.Context, symbol string, leverage int) error
	SetMarginMode(ctx context.Context, symbol string, marginMode MarginMode) error
	SetPositionMode(ctx context.Context, dualSidePosition bool) error

	// WebSocket interface
	SubscribeMarketData(ctx context.Context, symbols []string, dataTypes []string, callback EventCallback) error
	UnsubscribeMarketData(ctx context.Context, symbols []string, dataTypes []string) error
	SubscribeUserData(ctx context.Context, callback EventCallback) error
	UnsubscribeUserData(ctx context.Context) error

	// Status and monitoring
	GetStatus() ExchangeStatus
	GetConfig() *ExchangeConfig
	GetLastHeartbeat() time.Time
	GetMessageCount() int64
	GetErrorCount() int64

	// Event management
	RegisterEventCallback(eventType string, callback EventCallback)
	UnregisterEventCallback(eventType string)
}

// BaseExchange provides common functionality for exchange implementations
type BaseExchange struct {
	config             *ExchangeConfig
	status             ExchangeStatus
	eventCallbacks     map[string][]EventCallback
	lastHeartbeat      time.Time
	messageCount       int64
	errorCount         int64
	startTime          time.Time
}

// NewBaseExchange creates a new BaseExchange
func NewBaseExchange(config *ExchangeConfig) *BaseExchange {
	return &BaseExchange{
		config:         config,
		status:         ExchangeStatusDisconnected,
		eventCallbacks: make(map[string][]EventCallback),
		lastHeartbeat:  time.Now(),
		messageCount:   0,
		errorCount:     0,
		startTime:      time.Now(),
	}
}

// GetStatus returns the exchange status
func (e *BaseExchange) GetStatus() ExchangeStatus {
	return e.status
}

// SetStatus sets the exchange status
func (e *BaseExchange) SetStatus(status ExchangeStatus) {
	e.status = status
}

// GetConfig returns the exchange configuration
func (e *BaseExchange) GetConfig() *ExchangeConfig {
	return e.config
}

// GetLastHeartbeat returns the last heartbeat time
func (e *BaseExchange) GetLastHeartbeat() time.Time {
	return e.lastHeartbeat
}

// UpdateHeartbeat updates the last heartbeat time
func (e *BaseExchange) UpdateHeartbeat() {
	e.lastHeartbeat = time.Now()
}

// GetMessageCount returns the message count
func (e *BaseExchange) GetMessageCount() int64 {
	return e.messageCount
}

// IncrementMessageCount increments the message count
func (e *BaseExchange) IncrementMessageCount() {
	e.messageCount++
}

// GetErrorCount returns the error count
func (e *BaseExchange) GetErrorCount() int64 {
	return e.errorCount
}

// IncrementErrorCount increments the error count
func (e *BaseExchange) IncrementErrorCount() {
	e.errorCount++
}

// GetUptime returns the uptime in seconds
func (e *BaseExchange) GetUptime() float64 {
	return time.Since(e.startTime).Seconds()
}

// RegisterEventCallback registers an event callback
func (e *BaseExchange) RegisterEventCallback(eventType string, callback EventCallback) {
	if _, ok := e.eventCallbacks[eventType]; !ok {
		e.eventCallbacks[eventType] = []EventCallback{}
	}
	e.eventCallbacks[eventType] = append(e.eventCallbacks[eventType], callback)
}

// UnregisterEventCallback unregisters an event callback
func (e *BaseExchange) UnregisterEventCallback(eventType string) {
	delete(e.eventCallbacks, eventType)
}

// EmitEvent emits an event to registered callbacks
func (e *BaseExchange) EmitEvent(eventType string, event interface{}) {
	if callbacks, ok := e.eventCallbacks[eventType]; ok {
		for _, callback := range callbacks {
			go callback(event)
		}
	}
}
