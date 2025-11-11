package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	App       AppConfig       `mapstructure:"app"`
	Exchanges []ExchangeConfig `mapstructure:"exchanges"`
	Trading   TradingConfig   `mapstructure:"trading"`
	Logging   LoggingConfig   `mapstructure:"logging"`
}

// AppConfig represents application-level configuration
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
}

// ExchangeConfig represents exchange configuration
type ExchangeConfig struct {
	ID            string                 `mapstructure:"id"`
	Name          string                 `mapstructure:"name"`
	Type          string                 `mapstructure:"type"`
	APIKey        string                 `mapstructure:"api_key"`
	APISecret     string                 `mapstructure:"api_secret"`
	APIPassphrase string                 `mapstructure:"api_passphrase"`
	WalletAddress string                 `mapstructure:"wallet_address"`
	Testnet       bool                   `mapstructure:"testnet"`
	Enabled       bool                   `mapstructure:"enabled"`
	Extra         map[string]interface{} `mapstructure:"extra"`
}

// TradingConfig represents trading configuration
type TradingConfig struct {
	DefaultLeverage   int     `mapstructure:"default_leverage"`
	DefaultMarginMode string  `mapstructure:"default_margin_mode"`
	MaxPositionSize   float64 `mapstructure:"max_position_size"`
	RiskPerTrade      float64 `mapstructure:"risk_per_trade"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level      string `mapstructure:"level"`
	OutputPath string `mapstructure:"output_path"`
	Colored    bool   `mapstructure:"colored"`
}

// Manager manages configuration loading and access
type Manager struct {
	viper  *viper.Viper
	config *Config
}

// NewManager creates a new configuration manager
func NewManager() *Manager {
	v := viper.New()
	v.SetConfigType("yaml")
	v.AutomaticEnv()

	return &Manager{
		viper: v,
	}
}

// LoadFromFile loads configuration from a file
func (m *Manager) LoadFromFile(path string) error {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("config file does not exist: %s", path)
	}

	// Set config file
	m.viper.SetConfigFile(path)

	// Read config
	if err := m.viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal config
	var config Config
	if err := m.viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	m.config = &config
	return nil
}

// LoadFromDirectory loads configuration from a directory
func (m *Manager) LoadFromDirectory(dir string) error {
	// Look for config files in the directory
	patterns := []string{
		filepath.Join(dir, "config.yaml"),
		filepath.Join(dir, "config.yml"),
	}

	for _, pattern := range patterns {
		if _, err := os.Stat(pattern); err == nil {
			return m.LoadFromFile(pattern)
		}
	}

	return fmt.Errorf("no config file found in directory: %s", dir)
}

// GetConfig returns the loaded configuration
func (m *Manager) GetConfig() *Config {
	return m.config
}

// GetExchange returns exchange configuration by ID
func (m *Manager) GetExchange(id string) (*ExchangeConfig, error) {
	if m.config == nil {
		return nil, fmt.Errorf("config not loaded")
	}

	for _, exchange := range m.config.Exchanges {
		if exchange.ID == id {
			return &exchange, nil
		}
	}

	return nil, fmt.Errorf("exchange not found: %s", id)
}

// GetEnabledExchanges returns all enabled exchanges
func (m *Manager) GetEnabledExchanges() []ExchangeConfig {
	if m.config == nil {
		return []ExchangeConfig{}
	}

	var enabled []ExchangeConfig
	for _, exchange := range m.config.Exchanges {
		if exchange.Enabled {
			enabled = append(enabled, exchange)
		}
	}

	return enabled
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		App: AppConfig{
			Name:        "Crypto Trading System",
			Version:     "2.0.0-go",
			Environment: "development",
		},
		Trading: TradingConfig{
			DefaultLeverage:   1,
			DefaultMarginMode: "cross",
			MaxPositionSize:   10000,
			RiskPerTrade:      0.02,
		},
		Logging: LoggingConfig{
			Level:   "info",
			Colored: true,
		},
		Exchanges: []ExchangeConfig{},
	}
}
