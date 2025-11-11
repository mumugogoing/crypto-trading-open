package binance

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/mumugogoing/crypto-trading-open/pkg/adapters/exchanges"
	httpClient "github.com/mumugogoing/crypto-trading-open/pkg/infrastructure/http"
	"github.com/mumugogoing/crypto-trading-open/pkg/infrastructure/logging"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

const (
	// Binance API URLs
	BinanceMainnetBaseURL = "https://fapi.binance.com"
	BinanceTestnetBaseURL = "https://testnet.binancefuture.com"
	BinanceMainnetWsURL   = "wss://fstream.binance.com/ws"
	BinanceTestnetWsURL   = "wss://stream.binancefuture.com/ws"
)

// BinanceAdapter implements the Binance exchange adapter
type BinanceAdapter struct {
	*exchanges.BaseExchange
	logger     *zap.Logger
	httpClient *httpClient.Client
	signer     *httpClient.SignatureHelper
	baseURL    string
	wsURL      string
}

// NewBinanceAdapter creates a new Binance adapter
func NewBinanceAdapter(config *exchanges.ExchangeConfig) (*BinanceAdapter, error) {
	logger := logging.GetLogger().With(
		zap.String("exchange", "binance"),
		zap.String("exchange_id", config.ExchangeID),
	)

	// Determine URLs based on testnet flag
	baseURL := BinanceMainnetBaseURL
	wsURL := BinanceMainnetWsURL
	if config.Testnet {
		baseURL = BinanceTestnetBaseURL
		wsURL = BinanceTestnetWsURL
	}

	// Override with custom URLs if provided
	if config.BaseURL != "" {
		baseURL = config.BaseURL
	}
	if config.WsURL != "" {
		wsURL = config.WsURL
	}

	// Create HTTP client
	client := httpClient.NewClient(baseURL, time.Duration(config.RequestTimeout)*time.Second, logger)

	adapter := &BinanceAdapter{
		BaseExchange: exchanges.NewBaseExchange(config),
		logger:       logger,
		httpClient:   client,
		signer:       httpClient.NewSignatureHelper(),
		baseURL:      baseURL,
		wsURL:        wsURL,
	}

	return adapter, nil
}

// Connect connects to Binance
func (b *BinanceAdapter) Connect(ctx context.Context) error {
	b.SetStatus(exchanges.ExchangeStatusConnecting)
	b.logger.Info("Connecting to Binance...")

	// Test connectivity
	_, err := b.HealthCheck(ctx)
	if err != nil {
		b.SetStatus(exchanges.ExchangeStatusError)
		return fmt.Errorf("failed to connect to Binance: %w", err)
	}

	b.SetStatus(exchanges.ExchangeStatusConnected)
	b.logger.Info("Connected to Binance successfully")
	return nil
}

// Disconnect disconnects from Binance
func (b *BinanceAdapter) Disconnect(ctx context.Context) error {
	b.SetStatus(exchanges.ExchangeStatusDisconnected)
	b.logger.Info("Disconnected from Binance")
	return nil
}

// Authenticate authenticates with Binance
func (b *BinanceAdapter) Authenticate(ctx context.Context) error {
	// Test authentication by getting account info
	balance, err := b.GetBalance(ctx)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	b.SetStatus(exchanges.ExchangeStatusAuthenticated)
	b.logger.Info("Authenticated with Binance successfully", zap.Int("assets", len(balance)))
	return nil
}

// HealthCheck performs a health check
func (b *BinanceAdapter) HealthCheck(ctx context.Context) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := b.httpClient.Get(ctx, "/fapi/v1/ping", nil, &result)
	if err != nil {
		return nil, err
	}

	// Get server time
	var timeResult map[string]interface{}
	err = b.httpClient.Get(ctx, "/fapi/v1/time", nil, &timeResult)
	if err != nil {
		return nil, err
	}

	health := map[string]interface{}{
		"status":      "ok",
		"server_time": timeResult["serverTime"],
		"exchange":    "binance",
	}

	return health, nil
}

// GetExchangeInfo gets exchange information
func (b *BinanceAdapter) GetExchangeInfo(ctx context.Context) (*exchanges.ExchangeInfo, error) {
	var result map[string]interface{}
	err := b.httpClient.Get(ctx, "/fapi/v1/exchangeInfo", nil, &result)
	if err != nil {
		return nil, err
	}

	info := exchanges.NewExchangeInfo("Binance")
	
	// Extract symbols
	if symbols, ok := result["symbols"].([]interface{}); ok {
		for _, s := range symbols {
			if symbolMap, ok := s.(map[string]interface{}); ok {
				if symbol, ok := symbolMap["symbol"].(string); ok {
					info.Symbols = append(info.Symbols, symbol)
				}
			}
		}
	}

	info.RawData = result
	return info, nil
}

// GetTicker gets ticker data for a symbol
func (b *BinanceAdapter) GetTicker(ctx context.Context, symbol string) (*exchanges.TickerData, error) {
	params := map[string]string{
		"symbol": symbol,
	}

	var result map[string]interface{}
	err := b.httpClient.Get(ctx, "/fapi/v1/ticker/24hr", params, &result)
	if err != nil {
		return nil, err
	}

	ticker := exchanges.NewTickerData(symbol)
	
	if lastPrice, ok := result["lastPrice"].(string); ok {
		ticker.Last, _ = decimal.NewFromString(lastPrice)
	}
	if bidPrice, ok := result["bidPrice"].(string); ok {
		ticker.Bid, _ = decimal.NewFromString(bidPrice)
	}
	if askPrice, ok := result["askPrice"].(string); ok {
		ticker.Ask, _ = decimal.NewFromString(askPrice)
	}
	if highPrice, ok := result["highPrice"].(string); ok {
		ticker.High, _ = decimal.NewFromString(highPrice)
	}
	if lowPrice, ok := result["lowPrice"].(string); ok {
		ticker.Low, _ = decimal.NewFromString(lowPrice)
	}
	if volume, ok := result["volume"].(string); ok {
		ticker.Volume, _ = decimal.NewFromString(volume)
	}

	ticker.RawData = result
	return ticker, nil
}

// GetOrderBook gets orderbook data for a symbol
func (b *BinanceAdapter) GetOrderBook(ctx context.Context, symbol string, depth int) (*exchanges.OrderBookData, error) {
	params := map[string]string{
		"symbol": symbol,
		"limit":  strconv.Itoa(depth),
	}

	var result map[string]interface{}
	err := b.httpClient.Get(ctx, "/fapi/v1/depth", params, &result)
	if err != nil {
		return nil, err
	}

	orderbook := exchanges.NewOrderBookData(symbol)
	
	// Parse bids
	if bids, ok := result["bids"].([]interface{}); ok {
		for _, bid := range bids {
			if bidArray, ok := bid.([]interface{}); ok && len(bidArray) >= 2 {
				price, _ := decimal.NewFromString(bidArray[0].(string))
				amount, _ := decimal.NewFromString(bidArray[1].(string))
				orderbook.Bids = append(orderbook.Bids, []decimal.Decimal{price, amount})
			}
		}
	}

	// Parse asks
	if asks, ok := result["asks"].([]interface{}); ok {
		for _, ask := range asks {
			if askArray, ok := ask.([]interface{}); ok && len(askArray) >= 2 {
				price, _ := decimal.NewFromString(askArray[0].(string))
				amount, _ := decimal.NewFromString(askArray[1].(string))
				orderbook.Asks = append(orderbook.Asks, []decimal.Decimal{price, amount})
			}
		}
	}

	orderbook.RawData = result
	return orderbook, nil
}

// GetTrades gets recent trades for a symbol
func (b *BinanceAdapter) GetTrades(ctx context.Context, symbol string, limit int) ([]*exchanges.TradeData, error) {
	params := map[string]string{
		"symbol": symbol,
		"limit":  strconv.Itoa(limit),
	}

	var result []map[string]interface{}
	err := b.httpClient.Get(ctx, "/fapi/v1/trades", params, &result)
	if err != nil {
		return nil, err
	}

	trades := make([]*exchanges.TradeData, 0, len(result))
	for _, t := range result {
		trade := &exchanges.TradeData{
			Symbol:  symbol,
			RawData: t,
		}

		if id, ok := t["id"].(float64); ok {
			trade.ID = strconv.FormatFloat(id, 'f', 0, 64)
		}
		if price, ok := t["price"].(string); ok {
			trade.Price, _ = decimal.NewFromString(price)
		}
		if qty, ok := t["qty"].(string); ok {
			trade.Amount, _ = decimal.NewFromString(qty)
		}
		if isBuyerMaker, ok := t["isBuyerMaker"].(bool); ok {
			if isBuyerMaker {
				trade.Side = exchanges.OrderSideSell
			} else {
				trade.Side = exchanges.OrderSideBuy
			}
		}
		if tradeTime, ok := t["time"].(float64); ok {
			trade.Timestamp = time.Unix(int64(tradeTime)/1000, 0)
		}

		trades = append(trades, trade)
	}

	return trades, nil
}

// GetOHLCV gets OHLCV data for a symbol
func (b *BinanceAdapter) GetOHLCV(ctx context.Context, symbol, timeframe string, since time.Time, limit int) ([]*exchanges.OHLCVData, error) {
	params := map[string]string{
		"symbol":   symbol,
		"interval": timeframe,
		"limit":    strconv.Itoa(limit),
	}

	if !since.IsZero() {
		params["startTime"] = strconv.FormatInt(since.UnixMilli(), 10)
	}

	var result [][]interface{}
	err := b.httpClient.Get(ctx, "/fapi/v1/klines", params, &result)
	if err != nil {
		return nil, err
	}

	ohlcvs := make([]*exchanges.OHLCVData, 0, len(result))
	for _, k := range result {
		if len(k) < 6 {
			continue
		}

		timestamp := time.Unix(int64(k[0].(float64))/1000, 0)
		ohlcv := exchanges.NewOHLCVData(symbol, timestamp)
		
		if open, ok := k[1].(string); ok {
			ohlcv.Open, _ = decimal.NewFromString(open)
		}
		if high, ok := k[2].(string); ok {
			ohlcv.High, _ = decimal.NewFromString(high)
		}
		if low, ok := k[3].(string); ok {
			ohlcv.Low, _ = decimal.NewFromString(low)
		}
		if close, ok := k[4].(string); ok {
			ohlcv.Close, _ = decimal.NewFromString(close)
		}
		if volume, ok := k[5].(string); ok {
			ohlcv.Volume, _ = decimal.NewFromString(volume)
		}

		ohlcvs = append(ohlcvs, ohlcv)
	}

	return ohlcvs, nil
}

// signRequest signs a request with API key and secret
func (b *BinanceAdapter) signRequest(params map[string]string) map[string]string {
	if params == nil {
		params = make(map[string]string)
	}

	// Add timestamp
	params["timestamp"] = strconv.FormatInt(time.Now().UnixMilli(), 10)

	// Create query string
	queryString := ""
	for k, v := range params {
		if queryString != "" {
			queryString += "&"
		}
		queryString += fmt.Sprintf("%s=%s", k, v)
	}

	// Sign the query string
	signature := b.signer.SignHMACSHA256(queryString, b.GetConfig().APISecret)
	params["signature"] = signature

	// Add API key header
	b.httpClient.SetHeader("X-MBX-APIKEY", b.GetConfig().APIKey)

	return params
}
