package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// ExchangeData represents exchange data
type ExchangeData struct {
	ExchangeID  string
	Name        string
	BaseURL     string
	WsURL       string
	Testnet     bool
	Connected   bool
	LastUpdate  time.Time
}

// NewExchangeData creates a new ExchangeData
func NewExchangeData(exchangeID, name, baseURL, wsURL string, testnet bool) *ExchangeData {
	return &ExchangeData{
		ExchangeID: exchangeID,
		Name:       name,
		BaseURL:    baseURL,
		WsURL:      wsURL,
		Testnet:    testnet,
		Connected:  false,
		LastUpdate: time.Now(),
	}
}

// PriceData represents price data
type PriceData struct {
	Symbol     string
	Exchange   string
	Price      float64
	Volume     float64
	Timestamp  time.Time
	LastUpdate time.Time
}

// NewPriceData creates a new PriceData
func NewPriceData(symbol, exchange string, price, volume float64) *PriceData {
	now := time.Now()
	return &PriceData{
		Symbol:     symbol,
		Exchange:   exchange,
		Price:      price,
		Volume:     volume,
		Timestamp:  now,
		LastUpdate: now,
	}
}

// SpreadData represents spread data
type SpreadData struct {
	Symbol     string
	Exchange1  string
	Exchange2  string
	Price1     float64
	Price2     float64
	Spread     float64
	SpreadPct  float64
	Volume1    float64
	Volume2    float64
	Timestamp  time.Time
}

// NewSpreadData creates a new SpreadData
func NewSpreadData(symbol, exchange1, exchange2 string, price1, price2, volume1, volume2 float64) *SpreadData {
	spread := price2 - price1
	spreadPct := 0.0
	if price1 > 0 {
		spreadPct = (spread / price1) * 100
	}
	
	return &SpreadData{
		Symbol:    symbol,
		Exchange1: exchange1,
		Exchange2: exchange2,
		Price1:    price1,
		Price2:    price2,
		Spread:    spread,
		SpreadPct: spreadPct,
		Volume1:   volume1,
		Volume2:   volume2,
		Timestamp: time.Now(),
	}
}

// SymbolInfo represents trading pair information
type SymbolInfo struct {
	Symbol            string
	BaseCurrency      string
	QuoteCurrency     string
	ContractType      string
	PricePrecision    int
	QuantityPrecision int
	MinQuantity       decimal.Decimal
	MaxQuantity       decimal.Decimal
	MinPrice          decimal.Decimal
	MaxPrice          decimal.Decimal
	Active            bool
}

// NewSymbolInfo creates a new SymbolInfo
func NewSymbolInfo(symbol, base, quote, contractType string) *SymbolInfo {
	return &SymbolInfo{
		Symbol:            symbol,
		BaseCurrency:      base,
		QuoteCurrency:     quote,
		ContractType:      contractType,
		PricePrecision:    8,
		QuantityPrecision: 8,
		MinQuantity:       decimal.Zero,
		MaxQuantity:       decimal.Zero,
		MinPrice:          decimal.Zero,
		MaxPrice:          decimal.Zero,
		Active:            true,
	}
}

// MarketData represents market data
type MarketData struct {
	Symbol     string
	Exchange   string
	Ticker     map[string]interface{}
	Orderbook  map[string]interface{}
	Trades     []map[string]interface{}
	LastUpdate time.Time
}

// NewMarketData creates a new MarketData
func NewMarketData(symbol, exchange string) *MarketData {
	return &MarketData{
		Symbol:     symbol,
		Exchange:   exchange,
		Ticker:     make(map[string]interface{}),
		Orderbook:  make(map[string]interface{}),
		Trades:     []map[string]interface{}{},
		LastUpdate: time.Now(),
	}
}

// ExchangeStatus represents exchange status
type ExchangeStatus struct {
	ExchangeID         string
	Connected          bool
	Authenticated      bool
	WebsocketConnected bool
	LastHeartbeat      time.Time
	MessageCount       int64
	ErrorCount         int64
	Uptime             float64
}

// NewExchangeStatus creates a new ExchangeStatus
func NewExchangeStatus(exchangeID string) *ExchangeStatus {
	return &ExchangeStatus{
		ExchangeID:         exchangeID,
		Connected:          false,
		Authenticated:      false,
		WebsocketConnected: false,
		LastHeartbeat:      time.Now(),
		MessageCount:       0,
		ErrorCount:         0,
		Uptime:             0,
	}
}
