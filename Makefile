# Makefile for Crypto Trading System (Go)

# Variables
BINARY_DIR=bin
GO=go
GOFLAGS=-v
LDFLAGS=-s -w

# Binaries
TEST_BINANCE=$(BINARY_DIR)/test_binance
GRID_TRADING=$(BINARY_DIR)/grid_trading
VOLUME_MAKER=$(BINARY_DIR)/volume_maker
ARBITRAGE_MONITOR=$(BINARY_DIR)/arbitrage_monitor
PRICE_ALERT=$(BINARY_DIR)/price_alert

# Phony targets
.PHONY: all clean build test fmt vet lint help deps

# Default target
all: deps build

# Help target
help:
	@echo "Available targets:"
	@echo "  all              - Download dependencies and build all binaries"
	@echo "  build            - Build all binaries"
	@echo "  test             - Run tests"
	@echo "  test-binance     - Build and run Binance test"
	@echo "  clean            - Remove built binaries"
	@echo "  deps             - Download dependencies"
	@echo "  fmt              - Format code"
	@echo "  vet              - Run go vet"
	@echo "  lint             - Run golangci-lint (if installed)"
	@echo "  help             - Show this help message"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy

# Build all binaries
build: $(TEST_BINANCE)

# Build test_binance
$(TEST_BINANCE):
	@echo "Building test_binance..."
	@mkdir -p $(BINARY_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(TEST_BINANCE) ./cmd/test_binance/

# Build grid_trading (TODO: implement)
$(GRID_TRADING):
	@echo "Building grid_trading..."
	@mkdir -p $(BINARY_DIR)
	@echo "⚠️  grid_trading not yet implemented"
	# $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(GRID_TRADING) ./cmd/grid_trading/

# Build volume_maker (TODO: implement)
$(VOLUME_MAKER):
	@echo "Building volume_maker..."
	@mkdir -p $(BINARY_DIR)
	@echo "⚠️  volume_maker not yet implemented"
	# $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(VOLUME_MAKER) ./cmd/volume_maker/

# Build arbitrage_monitor (TODO: implement)
$(ARBITRAGE_MONITOR):
	@echo "Building arbitrage_monitor..."
	@mkdir -p $(BINARY_DIR)
	@echo "⚠️  arbitrage_monitor not yet implemented"
	# $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(ARBITRAGE_MONITOR) ./cmd/arbitrage_monitor/

# Build price_alert (TODO: implement)
$(PRICE_ALERT):
	@echo "Building price_alert..."
	@mkdir -p $(BINARY_DIR)
	@echo "⚠️  price_alert not yet implemented"
	# $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(PRICE_ALERT) ./cmd/price_alert/

# Run tests
test:
	@echo "Running tests..."
	$(GO) test -v -race -cover ./...

# Run Binance test
test-binance: $(TEST_BINANCE)
	@echo "Running Binance test..."
	./$(TEST_BINANCE)

# Clean built binaries
clean:
	@echo "Cleaning..."
	rm -rf $(BINARY_DIR)
	$(GO) clean

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	$(GO) vet ./...

# Run golangci-lint
lint:
	@echo "Running golangci-lint..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install from https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run

# Install tools
install-tools:
	@echo "Installing development tools..."
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
