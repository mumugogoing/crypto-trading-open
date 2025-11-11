package binance

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/mumugogoing/crypto-trading-open/pkg/adapters/exchanges"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// CreateOrder creates a new order
func (b *BinanceAdapter) CreateOrder(ctx context.Context, symbol string, side exchanges.OrderSide, orderType exchanges.OrderType, amount, price decimal.Decimal, params map[string]interface{}) (*exchanges.OrderData, error) {
	orderParams := map[string]string{
		"symbol":   symbol,
		"side":     string(side),
		"type":     string(orderType),
		"quantity": amount.String(),
	}

	// Add price for limit orders
	if orderType == exchanges.OrderTypeLimit || orderType == exchanges.OrderTypeStopLimit {
		orderParams["price"] = price.String()
	}

	// Add optional parameters
	if timeInForce, ok := params["timeInForce"].(string); ok {
		orderParams["timeInForce"] = timeInForce
	} else if orderType == exchanges.OrderTypeLimit {
		orderParams["timeInForce"] = "GTC"
	}

	// Sign the request
	orderParams = b.signRequest(orderParams)

	var result map[string]interface{}
	err := b.httpClient.Post(ctx, "/fapi/v1/order", orderParams, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return b.parseOrderData(result), nil
}

// CancelOrder cancels an order
func (b *BinanceAdapter) CancelOrder(ctx context.Context, orderID, symbol string) (*exchanges.OrderData, error) {
	params := map[string]string{
		"symbol":  symbol,
		"orderId": orderID,
	}

	params = b.signRequest(params)

	var result map[string]interface{}
	err := b.httpClient.Delete(ctx, "/fapi/v1/order", params, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel order: %w", err)
	}

	return b.parseOrderData(result), nil
}

// CancelAllOrders cancels all orders for a symbol
func (b *BinanceAdapter) CancelAllOrders(ctx context.Context, symbol string) ([]*exchanges.OrderData, error) {
	params := map[string]string{
		"symbol": symbol,
	}

	params = b.signRequest(params)

	var result map[string]interface{}
	err := b.httpClient.Delete(ctx, "/fapi/v1/allOpenOrders", params, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel all orders: %w", err)
	}

	b.logger.Info("Cancelled all orders", zap.String("symbol", symbol))
	return []*exchanges.OrderData{}, nil
}

// GetOrder gets an order by ID
func (b *BinanceAdapter) GetOrder(ctx context.Context, orderID, symbol string) (*exchanges.OrderData, error) {
	params := map[string]string{
		"symbol":  symbol,
		"orderId": orderID,
	}

	params = b.signRequest(params)

	var result map[string]interface{}
	err := b.httpClient.Get(ctx, "/fapi/v1/order", params, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return b.parseOrderData(result), nil
}

// GetOpenOrders gets all open orders
func (b *BinanceAdapter) GetOpenOrders(ctx context.Context, symbol string) ([]*exchanges.OrderData, error) {
	params := map[string]string{}
	if symbol != "" {
		params["symbol"] = symbol
	}

	params = b.signRequest(params)

	var result []map[string]interface{}
	err := b.httpClient.Get(ctx, "/fapi/v1/openOrders", params, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get open orders: %w", err)
	}

	orders := make([]*exchanges.OrderData, 0, len(result))
	for _, o := range result {
		orders = append(orders, b.parseOrderData(o))
	}

	return orders, nil
}

// GetClosedOrders gets closed orders
func (b *BinanceAdapter) GetClosedOrders(ctx context.Context, symbol string, since time.Time, limit int) ([]*exchanges.OrderData, error) {
	params := map[string]string{
		"symbol": symbol,
	}

	if limit > 0 {
		params["limit"] = strconv.Itoa(limit)
	}
	if !since.IsZero() {
		params["startTime"] = strconv.FormatInt(since.UnixMilli(), 10)
	}

	params = b.signRequest(params)

	var result []map[string]interface{}
	err := b.httpClient.Get(ctx, "/fapi/v1/allOrders", params, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get closed orders: %w", err)
	}

	orders := make([]*exchanges.OrderData, 0, len(result))
	for _, o := range result {
		orderData := b.parseOrderData(o)
		if orderData.Status == exchanges.OrderStatusFilled || orderData.Status == exchanges.OrderStatusCanceled {
			orders = append(orders, orderData)
		}
	}

	return orders, nil
}

// GetBalance gets account balance
func (b *BinanceAdapter) GetBalance(ctx context.Context) (map[string]*exchanges.BalanceData, error) {
	params := b.signRequest(nil)

	var result map[string]interface{}
	err := b.httpClient.Get(ctx, "/fapi/v2/balance", params, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	balances := make(map[string]*exchanges.BalanceData)
	
	if assets, ok := result["assets"].([]interface{}); ok {
		for _, asset := range assets {
			if assetMap, ok := asset.(map[string]interface{}); ok {
				currency := assetMap["asset"].(string)
				free, _ := decimal.NewFromString(assetMap["availableBalance"].(string))
				total, _ := decimal.NewFromString(assetMap["balance"].(string))
				used := total.Sub(free)

				balances[currency] = exchanges.NewBalanceData(currency, free, used, total)
			}
		}
	}

	return balances, nil
}

// GetPositions gets all positions
func (b *BinanceAdapter) GetPositions(ctx context.Context, symbol string) ([]*exchanges.PositionData, error) {
	params := b.signRequest(nil)

	var result []map[string]interface{}
	err := b.httpClient.Get(ctx, "/fapi/v2/positionRisk", params, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get positions: %w", err)
	}

	positions := make([]*exchanges.PositionData, 0)
	for _, p := range result {
		posSymbol := p["symbol"].(string)
		
		// Filter by symbol if specified
		if symbol != "" && posSymbol != symbol {
			continue
		}

		positionAmt, _ := decimal.NewFromString(p["positionAmt"].(string))
		
		// Skip zero positions
		if positionAmt.IsZero() {
			continue
		}

		entryPrice, _ := decimal.NewFromString(p["entryPrice"].(string))
		markPrice, _ := decimal.NewFromString(p["markPrice"].(string))
		unrealizedProfit, _ := decimal.NewFromString(p["unRealizedProfit"].(string))
		
		side := exchanges.PositionSideLong
		if positionAmt.IsNegative() {
			side = exchanges.PositionSideShort
			positionAmt = positionAmt.Abs()
		}

		position := exchanges.NewPositionData(posSymbol, side, positionAmt, entryPrice)
		position.MarkPrice = markPrice
		position.UnrealizedPnl = unrealizedProfit
		position.RawData = p

		if leverage, ok := p["leverage"].(string); ok {
			if lev, err := strconv.Atoi(leverage); err == nil {
				position.Leverage = lev
			}
		}

		positions = append(positions, position)
	}

	return positions, nil
}

// GetPosition gets a single position
func (b *BinanceAdapter) GetPosition(ctx context.Context, symbol string) (*exchanges.PositionData, error) {
	positions, err := b.GetPositions(ctx, symbol)
	if err != nil {
		return nil, err
	}

	if len(positions) == 0 {
		return nil, fmt.Errorf("position not found for symbol: %s", symbol)
	}

	return positions[0], nil
}

// SetLeverage sets leverage for a symbol
func (b *BinanceAdapter) SetLeverage(ctx context.Context, symbol string, leverage int) error {
	params := map[string]string{
		"symbol":   symbol,
		"leverage": strconv.Itoa(leverage),
	}

	params = b.signRequest(params)

	var result map[string]interface{}
	err := b.httpClient.Post(ctx, "/fapi/v1/leverage", params, &result)
	if err != nil {
		return fmt.Errorf("failed to set leverage: %w", err)
	}

	b.logger.Info("Set leverage", zap.String("symbol", symbol), zap.Int("leverage", leverage))
	return nil
}

// SetMarginMode sets margin mode for a symbol
func (b *BinanceAdapter) SetMarginMode(ctx context.Context, symbol string, marginMode exchanges.MarginMode) error {
	params := map[string]string{
		"symbol":     symbol,
		"marginType": string(marginMode),
	}

	params = b.signRequest(params)

	var result map[string]interface{}
	err := b.httpClient.Post(ctx, "/fapi/v1/marginType", params, &result)
	if err != nil {
		return fmt.Errorf("failed to set margin mode: %w", err)
	}

	b.logger.Info("Set margin mode", zap.String("symbol", symbol), zap.String("mode", string(marginMode)))
	return nil
}

// SetPositionMode sets position mode (one-way or hedge)
func (b *BinanceAdapter) SetPositionMode(ctx context.Context, dualSidePosition bool) error {
	params := map[string]string{
		"dualSidePosition": strconv.FormatBool(dualSidePosition),
	}

	params = b.signRequest(params)

	var result map[string]interface{}
	err := b.httpClient.Post(ctx, "/fapi/v1/positionSide/dual", params, &result)
	if err != nil {
		return fmt.Errorf("failed to set position mode: %w", err)
	}

	b.logger.Info("Set position mode", zap.Bool("dualSide", dualSidePosition))
	return nil
}

// parseOrderData parses order data from API response
func (b *BinanceAdapter) parseOrderData(data map[string]interface{}) *exchanges.OrderData {
	orderID := fmt.Sprintf("%v", data["orderId"])
	symbol := data["symbol"].(string)
	
	side := exchanges.OrderSideBuy
	if data["side"].(string) == "SELL" {
		side = exchanges.OrderSideSell
	}

	orderType := exchanges.OrderType(data["type"].(string))
	
	origQty, _ := decimal.NewFromString(data["origQty"].(string))
	executedQty, _ := decimal.NewFromString(data["executedQty"].(string))
	
	price := decimal.Zero
	if priceStr, ok := data["price"].(string); ok && priceStr != "" {
		price, _ = decimal.NewFromString(priceStr)
	}

	order := exchanges.NewOrderData(orderID, symbol, side, orderType, origQty, price)
	order.Filled = executedQty
	order.Remaining = origQty.Sub(executedQty)
	
	if avgPrice, ok := data["avgPrice"].(string); ok && avgPrice != "" {
		order.Average, _ = decimal.NewFromString(avgPrice)
	}

	status := data["status"].(string)
	switch status {
	case "NEW":
		order.Status = exchanges.OrderStatusOpen
	case "PARTIALLY_FILLED":
		order.Status = exchanges.OrderStatusOpen
	case "FILLED":
		order.Status = exchanges.OrderStatusFilled
	case "CANCELED":
		order.Status = exchanges.OrderStatusCanceled
	case "REJECTED":
		order.Status = exchanges.OrderStatusRejected
	case "EXPIRED":
		order.Status = exchanges.OrderStatusExpired
	default:
		order.Status = exchanges.OrderStatusUnknown
	}

	if updateTime, ok := data["updateTime"].(float64); ok {
		order.Updated = time.Unix(int64(updateTime)/1000, 0)
	}

	order.RawData = data
	return order
}
