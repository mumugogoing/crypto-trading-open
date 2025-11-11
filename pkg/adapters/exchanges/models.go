package exchanges

import (
	"time"

	"github.com/shopspring/decimal"
)

// ExchangeType represents exchange type
type ExchangeType string

const (
	ExchangeTypeSpot       ExchangeType = "spot"
	ExchangeTypeFutures    ExchangeType = "futures"
	ExchangeTypePerpetual  ExchangeType = "perpetual"
	ExchangeTypeOptions    ExchangeType = "options"
	ExchangeTypeMargin     ExchangeType = "margin"
)

// OrderSide represents order side
type OrderSide string

const (
	OrderSideBuy  OrderSide = "buy"
	OrderSideSell OrderSide = "sell"
)

// OrderType represents order type
type OrderType string

const (
	OrderTypeMarket          OrderType = "market"
	OrderTypeLimit           OrderType = "limit"
	OrderTypeStop            OrderType = "stop"
	OrderTypeStopLimit       OrderType = "stop_limit"
	OrderTypeTakeProfit      OrderType = "take_profit"
	OrderTypeTakeProfitLimit OrderType = "take_profit_limit"
	OrderTypeIOC             OrderType = "ioc"
	OrderTypeFOK             OrderType = "fok"
	OrderTypePostOnly        OrderType = "post_only"
)

// OrderStatus represents order status
type OrderStatus string

const (
	OrderStatusPending  OrderStatus = "pending"
	OrderStatusOpen     OrderStatus = "open"
	OrderStatusFilled   OrderStatus = "filled"
	OrderStatusCanceled OrderStatus = "canceled"
	OrderStatusRejected OrderStatus = "rejected"
	OrderStatusExpired  OrderStatus = "expired"
	OrderStatusUnknown  OrderStatus = "unknown"
)

// PositionSide represents position side
type PositionSide string

const (
	PositionSideLong  PositionSide = "long"
	PositionSideShort PositionSide = "short"
	PositionSideBoth  PositionSide = "both"
)

// MarginMode represents margin mode
type MarginMode string

const (
	MarginModeCross    MarginMode = "cross"
	MarginModeIsolated MarginMode = "isolated"
)

// OrderData represents order data
type OrderData struct {
	ID        string
	ClientID  string
	Symbol    string
	Side      OrderSide
	Type      OrderType
	Amount    decimal.Decimal
	Price     decimal.Decimal
	Filled    decimal.Decimal
	Remaining decimal.Decimal
	Cost      decimal.Decimal
	Average   decimal.Decimal
	Status    OrderStatus
	Timestamp time.Time
	Updated   time.Time
	Fee       map[string]interface{}
	Trades    []map[string]interface{}
	Params    map[string]interface{}
	RawData   map[string]interface{}
}

// NewOrderData creates a new OrderData
func NewOrderData(id, symbol string, side OrderSide, orderType OrderType, amount, price decimal.Decimal) *OrderData {
	now := time.Now()
	return &OrderData{
		ID:        id,
		Symbol:    symbol,
		Side:      side,
		Type:      orderType,
		Amount:    amount,
		Price:     price,
		Filled:    decimal.Zero,
		Remaining: amount,
		Cost:      decimal.Zero,
		Average:   decimal.Zero,
		Status:    OrderStatusPending,
		Timestamp: now,
		Updated:   now,
		Fee:       make(map[string]interface{}),
		Trades:    []map[string]interface{}{},
		Params:    make(map[string]interface{}),
		RawData:   make(map[string]interface{}),
	}
}

// PositionData represents position data
type PositionData struct {
	Symbol           string
	Side             PositionSide
	Size             decimal.Decimal
	EntryPrice       decimal.Decimal
	MarkPrice        decimal.Decimal
	CurrentPrice     decimal.Decimal
	UnrealizedPnl    decimal.Decimal
	RealizedPnl      decimal.Decimal
	Percentage       decimal.Decimal
	Leverage         int
	MarginMode       MarginMode
	Margin           decimal.Decimal
	LiquidationPrice decimal.Decimal
	Timestamp        time.Time
	RawData          map[string]interface{}
}

// NewPositionData creates a new PositionData
func NewPositionData(symbol string, side PositionSide, size, entryPrice decimal.Decimal) *PositionData {
	return &PositionData{
		Symbol:           symbol,
		Side:             side,
		Size:             size,
		EntryPrice:       entryPrice,
		MarkPrice:        decimal.Zero,
		CurrentPrice:     decimal.Zero,
		UnrealizedPnl:    decimal.Zero,
		RealizedPnl:      decimal.Zero,
		Percentage:       decimal.Zero,
		Leverage:         1,
		MarginMode:       MarginModeCross,
		Margin:           decimal.Zero,
		LiquidationPrice: decimal.Zero,
		Timestamp:        time.Now(),
		RawData:          make(map[string]interface{}),
	}
}

// BalanceData represents balance data
type BalanceData struct {
	Currency  string
	Free      decimal.Decimal
	Used      decimal.Decimal
	Total     decimal.Decimal
	Timestamp time.Time
	RawData   map[string]interface{}
}

// NewBalanceData creates a new BalanceData
func NewBalanceData(currency string, free, used, total decimal.Decimal) *BalanceData {
	return &BalanceData{
		Currency:  currency,
		Free:      free,
		Used:      used,
		Total:     total,
		Timestamp: time.Now(),
		RawData:   make(map[string]interface{}),
	}
}

// TickerData represents ticker data
type TickerData struct {
	Symbol    string
	Bid       decimal.Decimal
	Ask       decimal.Decimal
	Last      decimal.Decimal
	High      decimal.Decimal
	Low       decimal.Decimal
	Volume    decimal.Decimal
	Timestamp time.Time
	RawData   map[string]interface{}
}

// NewTickerData creates a new TickerData
func NewTickerData(symbol string) *TickerData {
	return &TickerData{
		Symbol:    symbol,
		Bid:       decimal.Zero,
		Ask:       decimal.Zero,
		Last:      decimal.Zero,
		High:      decimal.Zero,
		Low:       decimal.Zero,
		Volume:    decimal.Zero,
		Timestamp: time.Now(),
		RawData:   make(map[string]interface{}),
	}
}

// OHLCVData represents OHLCV data
type OHLCVData struct {
	Symbol    string
	Timestamp time.Time
	Open      decimal.Decimal
	High      decimal.Decimal
	Low       decimal.Decimal
	Close     decimal.Decimal
	Volume    decimal.Decimal
	RawData   map[string]interface{}
}

// NewOHLCVData creates a new OHLCVData
func NewOHLCVData(symbol string, timestamp time.Time) *OHLCVData {
	return &OHLCVData{
		Symbol:    symbol,
		Timestamp: timestamp,
		Open:      decimal.Zero,
		High:      decimal.Zero,
		Low:       decimal.Zero,
		Close:     decimal.Zero,
		Volume:    decimal.Zero,
		RawData:   make(map[string]interface{}),
	}
}

// OrderBookData represents orderbook data
type OrderBookData struct {
	Symbol    string
	Bids      [][]decimal.Decimal
	Asks      [][]decimal.Decimal
	Timestamp time.Time
	RawData   map[string]interface{}
}

// NewOrderBookData creates a new OrderBookData
func NewOrderBookData(symbol string) *OrderBookData {
	return &OrderBookData{
		Symbol:    symbol,
		Bids:      [][]decimal.Decimal{},
		Asks:      [][]decimal.Decimal{},
		Timestamp: time.Now(),
		RawData:   make(map[string]interface{}),
	}
}

// TradeData represents trade data
type TradeData struct {
	ID        string
	Symbol    string
	Side      OrderSide
	Price     decimal.Decimal
	Amount    decimal.Decimal
	Timestamp time.Time
	RawData   map[string]interface{}
}

// NewTradeData creates a new TradeData
func NewTradeData(id, symbol string, side OrderSide, price, amount decimal.Decimal) *TradeData {
	return &TradeData{
		ID:        id,
		Symbol:    symbol,
		Side:      side,
		Price:     price,
		Amount:    amount,
		Timestamp: time.Now(),
		RawData:   make(map[string]interface{}),
	}
}

// ExchangeInfo represents exchange information
type ExchangeInfo struct {
	Name      string
	Symbols   []string
	RateLimit map[string]int
	Timestamp time.Time
	RawData   map[string]interface{}
}

// NewExchangeInfo creates a new ExchangeInfo
func NewExchangeInfo(name string) *ExchangeInfo {
	return &ExchangeInfo{
		Name:      name,
		Symbols:   []string{},
		RateLimit: make(map[string]int),
		Timestamp: time.Now(),
		RawData:   make(map[string]interface{}),
	}
}
