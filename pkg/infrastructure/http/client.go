package http

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

// Client wraps resty client with additional functionality
type Client struct {
	client *resty.Client
	logger *zap.Logger
}

// NewClient creates a new HTTP client
func NewClient(baseURL string, timeout time.Duration, logger *zap.Logger) *Client {
	client := resty.New()
	client.SetBaseURL(baseURL)
	client.SetTimeout(timeout)
	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("User-Agent", "CryptoTradingSystem/2.0-Go")

	return &Client{
		client: client,
		logger: logger,
	}
}

// Get performs a GET request
func (c *Client) Get(ctx context.Context, url string, queryParams map[string]string, result interface{}) error {
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(queryParams).
		SetResult(result).
		Get(url)

	if err != nil {
		c.logger.Error("HTTP GET request failed", zap.String("url", url), zap.Error(err))
		return err
	}

	if resp.IsError() {
		c.logger.Error("HTTP GET response error",
			zap.String("url", url),
			zap.Int("status", resp.StatusCode()),
			zap.String("body", string(resp.Body())))
		return fmt.Errorf("HTTP error: %d - %s", resp.StatusCode(), string(resp.Body()))
	}

	return nil
}

// Post performs a POST request
func (c *Client) Post(ctx context.Context, url string, body interface{}, result interface{}) error {
	resp, err := c.client.R().
		SetContext(ctx).
		SetBody(body).
		SetResult(result).
		Post(url)

	if err != nil {
		c.logger.Error("HTTP POST request failed", zap.String("url", url), zap.Error(err))
		return err
	}

	if resp.IsError() {
		c.logger.Error("HTTP POST response error",
			zap.String("url", url),
			zap.Int("status", resp.StatusCode()),
			zap.String("body", string(resp.Body())))
		return fmt.Errorf("HTTP error: %d - %s", resp.StatusCode(), string(resp.Body()))
	}

	return nil
}

// Delete performs a DELETE request
func (c *Client) Delete(ctx context.Context, url string, queryParams map[string]string, result interface{}) error {
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(queryParams).
		SetResult(result).
		Delete(url)

	if err != nil {
		c.logger.Error("HTTP DELETE request failed", zap.String("url", url), zap.Error(err))
		return err
	}

	if resp.IsError() {
		c.logger.Error("HTTP DELETE response error",
			zap.String("url", url),
			zap.Int("status", resp.StatusCode()),
			zap.String("body", string(resp.Body())))
		return fmt.Errorf("HTTP error: %d - %s", resp.StatusCode(), string(resp.Body()))
	}

	return nil
}

// SetHeader sets a header for all requests
func (c *Client) SetHeader(key, value string) {
	c.client.SetHeader(key, value)
}

// SetHeaders sets multiple headers
func (c *Client) SetHeaders(headers map[string]string) {
	c.client.SetHeaders(headers)
}

// GetUnderlyingClient returns the underlying resty client
func (c *Client) GetUnderlyingClient() *resty.Client {
	return c.client
}

// SignatureHelper provides HMAC signing utilities
type SignatureHelper struct{}

// NewSignatureHelper creates a new signature helper
func NewSignatureHelper() *SignatureHelper {
	return &SignatureHelper{}
}

// SignHMACSHA256 signs a message using HMAC SHA256
func (h *SignatureHelper) SignHMACSHA256(message, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

// SignHMACSHA256Bytes signs bytes using HMAC SHA256
func (h *SignatureHelper) SignHMACSHA256Bytes(message []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(message)
	return hex.EncodeToString(mac.Sum(nil))
}

// RateLimiter provides rate limiting functionality
type RateLimiter struct {
	requestsPerSecond int
	ticker            *time.Ticker
	tokens            chan struct{}
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerSecond int) *RateLimiter {
	limiter := &RateLimiter{
		requestsPerSecond: requestsPerSecond,
		tokens:            make(chan struct{}, requestsPerSecond),
	}

	// Fill initial tokens
	for i := 0; i < requestsPerSecond; i++ {
		limiter.tokens <- struct{}{}
	}

	// Refill tokens
	limiter.ticker = time.NewTicker(time.Second / time.Duration(requestsPerSecond))
	go func() {
		for range limiter.ticker.C {
			select {
			case limiter.tokens <- struct{}{}:
			default:
			}
		}
	}()

	return limiter
}

// Wait waits for a token
func (r *RateLimiter) Wait() {
	<-r.tokens
}

// Stop stops the rate limiter
func (r *RateLimiter) Stop() {
	r.ticker.Stop()
}

// RequestMiddleware provides middleware for requests
func RequestMiddleware(logger *zap.Logger) resty.RequestMiddleware {
	return func(c *resty.Client, r *resty.Request) error {
		logger.Debug("HTTP Request",
			zap.String("method", r.Method),
			zap.String("url", r.URL),
			zap.Any("query", r.QueryParam),
		)
		return nil
	}
}

// ResponseMiddleware provides middleware for responses
func ResponseMiddleware(logger *zap.Logger) resty.ResponseMiddleware {
	return func(c *resty.Client, r *resty.Response) error {
		logger.Debug("HTTP Response",
			zap.String("url", r.Request.URL),
			zap.Int("status", r.StatusCode()),
			zap.Duration("time", r.Time()),
		)
		return nil
	}
}

// ErrorHook provides error handling
func ErrorHook(logger *zap.Logger) resty.ErrorHook {
	return func(req *resty.Request, err error) {
		if v, ok := err.(*resty.ResponseError); ok {
			logger.Error("HTTP Error",
				zap.String("url", req.URL),
				zap.Int("status", v.Response.StatusCode()),
				zap.Error(err),
			)
		}
	}
}

// RetryCondition determines if a request should be retried
func RetryCondition() resty.RetryConditionFunc {
	return func(r *resty.Response, err error) bool {
		// Retry on network errors
		if err != nil {
			return true
		}

		// Retry on 5xx errors and 429 (rate limit)
		statusCode := r.StatusCode()
		return statusCode >= 500 || statusCode == http.StatusTooManyRequests
	}
}
