package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mumugogoing/crypto-trading-open/pkg/adapters/exchanges"
	"github.com/mumugogoing/crypto-trading-open/pkg/adapters/exchanges/binance"
	"github.com/mumugogoing/crypto-trading-open/pkg/infrastructure/logging"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logConfig := &logging.Config{
		Level:   "info",
		Colored: true,
	}
	
	if err := logging.InitLogger(logConfig); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logging.Sync()

	logger := logging.GetLogger()
	logger.Info("Starting Binance adapter test...")

	// Create Binance configuration
	config := exchanges.NewExchangeConfig(
		"binance",
		"Binance",
		exchanges.ExchangeTypePerpetual,
		"",  // API key (empty for public data)
		"",  // API secret (empty for public data)
	)
	config.Testnet = false

	// Create adapter
	adapter, err := binance.NewBinanceAdapter(config)
	if err != nil {
		logger.Fatal("Failed to create adapter", zap.Error(err))
	}

	ctx := context.Background()

	// Connect to Binance
	logger.Info("Connecting to Binance...")
	if err := adapter.Connect(ctx); err != nil {
		logger.Fatal("Failed to connect", zap.Error(err))
	}

	// Test health check
	logger.Info("Testing health check...")
	health, err := adapter.HealthCheck(ctx)
	if err != nil {
		logger.Error("Health check failed", zap.Error(err))
	} else {
		logger.Info("Health check successful", zap.Any("health", health))
	}

	// Test get exchange info
	logger.Info("Getting exchange info...")
	info, err := adapter.GetExchangeInfo(ctx)
	if err != nil {
		logger.Error("Failed to get exchange info", zap.Error(err))
	} else {
		logger.Info("Exchange info retrieved",
			zap.String("name", info.Name),
			zap.Int("symbols_count", len(info.Symbols)))
		
		// Show first 5 symbols
		if len(info.Symbols) > 0 {
			logger.Info("Sample symbols",
				zap.Strings("symbols", info.Symbols[:min(5, len(info.Symbols))]))
		}
	}

	// Test get ticker
	symbol := "BTCUSDT"
	logger.Info("Getting ticker...", zap.String("symbol", symbol))
	ticker, err := adapter.GetTicker(ctx, symbol)
	if err != nil {
		logger.Error("Failed to get ticker", zap.Error(err))
	} else {
		logger.Info("Ticker retrieved",
			zap.String("symbol", ticker.Symbol),
			zap.String("last", ticker.Last.String()),
			zap.String("bid", ticker.Bid.String()),
			zap.String("ask", ticker.Ask.String()),
			zap.String("volume", ticker.Volume.String()))
	}

	// Test get orderbook
	logger.Info("Getting orderbook...", zap.String("symbol", symbol))
	orderbook, err := adapter.GetOrderBook(ctx, symbol, 10)
	if err != nil {
		logger.Error("Failed to get orderbook", zap.Error(err))
	} else {
		logger.Info("Orderbook retrieved",
			zap.String("symbol", orderbook.Symbol),
			zap.Int("bids_count", len(orderbook.Bids)),
			zap.Int("asks_count", len(orderbook.Asks)))
		
		if len(orderbook.Bids) > 0 {
			logger.Info("Best bid",
				zap.String("price", orderbook.Bids[0][0].String()),
				zap.String("amount", orderbook.Bids[0][1].String()))
		}
		if len(orderbook.Asks) > 0 {
			logger.Info("Best ask",
				zap.String("price", orderbook.Asks[0][0].String()),
				zap.String("amount", orderbook.Asks[0][1].String()))
		}
	}

	// Test get trades
	logger.Info("Getting recent trades...", zap.String("symbol", symbol))
	trades, err := adapter.GetTrades(ctx, symbol, 5)
	if err != nil {
		logger.Error("Failed to get trades", zap.Error(err))
	} else {
		logger.Info("Trades retrieved", zap.Int("count", len(trades)))
		for i, trade := range trades {
			logger.Info("Trade",
				zap.Int("index", i),
				zap.String("id", trade.ID),
				zap.String("side", string(trade.Side)),
				zap.String("price", trade.Price.String()),
				zap.String("amount", trade.Amount.String()))
		}
	}

	// Test get OHLCV
	logger.Info("Getting OHLCV data...", zap.String("symbol", symbol))
	ohlcvs, err := adapter.GetOHLCV(ctx, symbol, "1h", time.Now().Add(-24*time.Hour), 10)
	if err != nil {
		logger.Error("Failed to get OHLCV", zap.Error(err))
	} else {
		logger.Info("OHLCV data retrieved", zap.Int("count", len(ohlcvs)))
		if len(ohlcvs) > 0 {
			latest := ohlcvs[len(ohlcvs)-1]
			logger.Info("Latest candle",
				zap.Time("time", latest.Timestamp),
				zap.String("open", latest.Open.String()),
				zap.String("high", latest.High.String()),
				zap.String("low", latest.Low.String()),
				zap.String("close", latest.Close.String()),
				zap.String("volume", latest.Volume.String()))
		}
	}

	// Disconnect
	logger.Info("Disconnecting...")
	if err := adapter.Disconnect(ctx); err != nil {
		logger.Error("Failed to disconnect", zap.Error(err))
	}

	logger.Info("Test completed successfully!")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
