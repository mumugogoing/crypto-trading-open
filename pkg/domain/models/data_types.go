package models

import (
	"fmt"
	"strings"
	"time"
)

// DataType represents supported data types
type DataType string

const (
	DataTypeTicker    DataType = "ticker"
	DataTypeOrderbook DataType = "orderbook"
	DataTypeTrades    DataType = "trades"
	DataTypeUserData  DataType = "user_data"
)

// GetAllTypes returns all supported data types
func GetAllDataTypes() []DataType {
	return []DataType{
		DataTypeTicker,
		DataTypeOrderbook,
		DataTypeTrades,
		DataTypeUserData,
	}
}

// FromString creates DataType from string
func DataTypeFromString(value string) (DataType, error) {
	normalized := strings.ToLower(value)
	for _, dt := range GetAllDataTypes() {
		if string(dt) == normalized {
			return dt, nil
		}
	}
	return "", fmt.Errorf("unsupported data type: %s", value)
}

// DataTypeConfig represents data type configuration
type DataTypeConfig struct {
	EnabledTypes  map[DataType]bool
	DisabledTypes map[DataType]bool
}

// NewDataTypeConfig creates a new DataTypeConfig with defaults
func NewDataTypeConfig() *DataTypeConfig {
	config := &DataTypeConfig{
		EnabledTypes:  make(map[DataType]bool),
		DisabledTypes: make(map[DataType]bool),
	}
	
	// Default: enable ticker and orderbook
	if len(config.EnabledTypes) == 0 && len(config.DisabledTypes) == 0 {
		config.EnabledTypes[DataTypeTicker] = true
		config.EnabledTypes[DataTypeOrderbook] = true
	}
	
	return config
}

// IsEnabled checks if a data type is enabled
func (c *DataTypeConfig) IsEnabled(dataType DataType) bool {
	if c.DisabledTypes[dataType] {
		return false
	}
	if len(c.EnabledTypes) > 0 {
		return c.EnabledTypes[dataType]
	}
	return true
}

// GetEnabledTypes returns list of enabled data types
func (c *DataTypeConfig) GetEnabledTypes() []DataType {
	var types []DataType
	
	if len(c.EnabledTypes) > 0 {
		for dt := range c.EnabledTypes {
			if !c.DisabledTypes[dt] {
				types = append(types, dt)
			}
		}
	} else {
		// Return all types except disabled
		for _, dt := range GetAllDataTypes() {
			if !c.DisabledTypes[dt] {
				types = append(types, dt)
			}
		}
	}
	
	return types
}

// GetEnabledTypeNames returns list of enabled data type names
func (c *DataTypeConfig) GetEnabledTypeNames() []string {
	types := c.GetEnabledTypes()
	names := make([]string, len(types))
	for i, dt := range types {
		names[i] = string(dt)
	}
	return names
}

// ExchangeDataTypeConfig represents exchange-specific data type configuration
type ExchangeDataTypeConfig struct {
	ExchangeID         string
	DataTypes          *DataTypeConfig
	MaxSymbolsPerType  map[DataType]int
	PrioritySymbols    []string
}

// NewExchangeDataTypeConfig creates a new ExchangeDataTypeConfig
func NewExchangeDataTypeConfig(exchangeID string) *ExchangeDataTypeConfig {
	return &ExchangeDataTypeConfig{
		ExchangeID:        exchangeID,
		DataTypes:         NewDataTypeConfig(),
		MaxSymbolsPerType: make(map[DataType]int),
		PrioritySymbols:   []string{},
	}
}

// GetMaxSymbols returns max symbols for a data type
func (c *ExchangeDataTypeConfig) GetMaxSymbols(dataType DataType) (int, bool) {
	val, ok := c.MaxSymbolsPerType[dataType]
	return val, ok
}

// SetMaxSymbols sets max symbols for a data type
func (c *ExchangeDataTypeConfig) SetMaxSymbols(dataType DataType, maxSymbols int) {
	c.MaxSymbolsPerType[dataType] = maxSymbols
}

// MonitoringDataTypeConfig represents monitoring data type configuration
type MonitoringDataTypeConfig struct {
	GlobalEnabledTypes map[DataType]bool
	ExchangeConfigs    map[string]*ExchangeDataTypeConfig
}

// NewMonitoringDataTypeConfig creates a new MonitoringDataTypeConfig
func NewMonitoringDataTypeConfig() *MonitoringDataTypeConfig {
	return &MonitoringDataTypeConfig{
		GlobalEnabledTypes: map[DataType]bool{
			DataTypeTicker:    true,
			DataTypeOrderbook: true,
		},
		ExchangeConfigs: make(map[string]*ExchangeDataTypeConfig),
	}
}

// GetExchangeConfig returns exchange config
func (c *MonitoringDataTypeConfig) GetExchangeConfig(exchangeID string) (*ExchangeDataTypeConfig, bool) {
	config, ok := c.ExchangeConfigs[exchangeID]
	return config, ok
}

// SetExchangeConfig sets exchange config
func (c *MonitoringDataTypeConfig) SetExchangeConfig(exchangeID string, config *ExchangeDataTypeConfig) {
	c.ExchangeConfigs[exchangeID] = config
}

// GetEnabledTypesForExchange returns enabled data types for an exchange
func (c *MonitoringDataTypeConfig) GetEnabledTypesForExchange(exchangeID string) []DataType {
	if config, ok := c.GetExchangeConfig(exchangeID); ok {
		return config.DataTypes.GetEnabledTypes()
	}
	
	// Use global config
	var types []DataType
	for dt, enabled := range c.GlobalEnabledTypes {
		if enabled {
			types = append(types, dt)
		}
	}
	return types
}

// SubscriptionStatus represents subscription status
type SubscriptionStatus struct {
	ExchangeID   string
	Symbol       string
	DataType     DataType
	Status       string // pending, active, error, cancelled
	ErrorMessage string
	LastUpdate   time.Time
}

// NewSubscriptionStatus creates a new SubscriptionStatus
func NewSubscriptionStatus(exchangeID, symbol string, dataType DataType) *SubscriptionStatus {
	return &SubscriptionStatus{
		ExchangeID: exchangeID,
		Symbol:     symbol,
		DataType:   dataType,
		Status:     "pending",
		LastUpdate: time.Now(),
	}
}

// IsActive checks if subscription is active
func (s *SubscriptionStatus) IsActive() bool {
	return s.Status == "active"
}

// IsError checks if subscription has error
func (s *SubscriptionStatus) IsError() bool {
	return s.Status == "error"
}

// SubscriptionSummary represents subscription summary
type SubscriptionSummary struct {
	TotalSubscriptions   int
	ActiveSubscriptions  int
	ErrorSubscriptions   int
	PendingSubscriptions int
	ByExchange           map[string]map[string]int
	ByDataType           map[DataType]map[string]int
}

// NewSubscriptionSummary creates a new SubscriptionSummary
func NewSubscriptionSummary() *SubscriptionSummary {
	return &SubscriptionSummary{
		ByExchange: make(map[string]map[string]int),
		ByDataType: make(map[DataType]map[string]int),
	}
}

// UpdateFromStatus updates summary from subscription status
func (s *SubscriptionSummary) UpdateFromStatus(status *SubscriptionStatus) {
	s.TotalSubscriptions++
	
	if status.IsActive() {
		s.ActiveSubscriptions++
	} else if status.IsError() {
		s.ErrorSubscriptions++
	} else {
		s.PendingSubscriptions++
	}
	
	// Update by exchange
	if _, ok := s.ByExchange[status.ExchangeID]; !ok {
		s.ByExchange[status.ExchangeID] = map[string]int{
			"total": 0, "active": 0, "error": 0, "pending": 0,
		}
	}
	s.ByExchange[status.ExchangeID]["total"]++
	if status.IsActive() {
		s.ByExchange[status.ExchangeID]["active"]++
	} else if status.IsError() {
		s.ByExchange[status.ExchangeID]["error"]++
	} else {
		s.ByExchange[status.ExchangeID]["pending"]++
	}
	
	// Update by data type
	if _, ok := s.ByDataType[status.DataType]; !ok {
		s.ByDataType[status.DataType] = map[string]int{
			"total": 0, "active": 0, "error": 0, "pending": 0,
		}
	}
	s.ByDataType[status.DataType]["total"]++
	if status.IsActive() {
		s.ByDataType[status.DataType]["active"]++
	} else if status.IsError() {
		s.ByDataType[status.DataType]["error"]++
	} else {
		s.ByDataType[status.DataType]["pending"]++
	}
}
