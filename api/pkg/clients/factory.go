package clients

import (
	"fmt"
	"strings"
)

// ClientFactory creates StockDataClients based on the name of the provider. Contains a map of configurations, (ProviderName: ClientConfig).
type ClientFactory struct {
	config map[string]ClientConfig
}

// NewClientFactory Creates a new instance of the ClientFactory struct.
func NewClientFactory() *ClientFactory {
	return &ClientFactory{
		config: make(map[string]ClientConfig),
	}
}

// RegisterProvider Adds a new provider configuration to the factory.
func (cf *ClientFactory) RegisterProvider(name, baseURL, apiKey string) {
	cf.config[strings.ToLower(name)] = ClientConfig{
		BaseURL: baseURL,
		APIKey:  apiKey,
	}
}

// CreateClient creates a client for the specified provider. Return an error if the provider is not registered or implemented.
func (cf *ClientFactory) CreateClient(providerName string) (StockDataClient, error) {
	config, exists := cf.config[strings.ToLower(providerName)]
	if !exists {
		return nil, fmt.Errorf("Unknown provider: %s", providerName)
	}

	// Create the appropriate client implementation
	switch strings.ToLower(providerName) {
	case "alphavantage":
		return NewAlphaVantageClient(config.BaseURL, config.APIKey), nil
	case "fmp":
		return NewFMPClient(config.BaseURL, config.APIKey), nil
	default:
		return nil, fmt.Errorf("Provider %s not implemented", providerName)
	}
}

// GetRegisteredProviders returns a list of all registered provider names.
func (cf *ClientFactory) GetRegisteredProviders() []string {
	providers := make([]string, 0, len(cf.config))
	for name := range cf.config {
		providers = append(providers, name)
	}
	return providers
}
