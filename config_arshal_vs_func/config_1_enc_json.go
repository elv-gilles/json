package config_arshal_vs_func

import (
	"encoding/json"
	"fmt"
)

type AppConfig1 struct {
	QSpacesConfig1
}

// InitDefaults initializes the default configuration.
func (c *AppConfig1) InitDefaults() *AppConfig1 {
	c.QSpacesConfig1.InitDefaults()
	return c
}

func (c *AppConfig1) UnmarshalJSON(bts []byte) error {
	c.InitDefaults()

	type alias AppConfig1
	return json.Unmarshal(bts, (*alias)(c))
}

type QSpacesConfig1 struct {
	QSpaces []QSpaceConfig1 `json:"qspaces"`
}

func (c *QSpacesConfig1) Validate() error {
	names := make(map[string]bool)
	for _, sp := range c.QSpaces {
		if len(sp.Names) > 0 {
			for _, n := range sp.Names {
				_, ok := names[n]
				if ok {
					return fmt.Errorf("duplicate space name %s", n)
				}
				names[n] = true
			}
		}
		if err := sp.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (c *QSpacesConfig1) InitDefaults() *QSpacesConfig1 {
	c.QSpaces = make([]QSpaceConfig1, 1)
	c.QSpaces[0] = QSpaceConfig1{
		ID:   "ispcMockContentSpace",
		Type: "Mock",
		Mock: (&MockQSpaceConfig1{}).InitDefaults(),
	}
	c.QSpaces[0].Services.InitDefaults()
	return c
}

type QSpaceConfig1 struct {
	ID       string                 `json:"id"`             // space ID
	Names    []string               `json:"names"`          // space names
	Type     string                 `json:"type"`           // Type is one of the supported content space types
	Mock     *MockQSpaceConfig1     `json:"mock,omitempty"` // Mock is the configuration for a mock content space
	Ethereum *EthereumQSpaceConfig1 `json:"ethereum,omitempty"`
	Services SpaceServicesConfig1   `json:"services"`
}

func (c *QSpaceConfig1) InitDefaults() *QSpaceConfig1 {
	c.Services.InitDefaults()
	return c
}

func (c *QSpaceConfig1) Validate() error {
	if c.Type == "Ethereum" {
		if c.Ethereum == nil {
			return fmt.Errorf("ethereum not defined")
		}
		return c.Ethereum.Validate()
	}
	return nil
}

func (c *QSpaceConfig1) UnmarshalJSON(bts []byte) error {
	c.InitDefaults()

	type alias QSpaceConfig1
	return json.Unmarshal(bts, (*alias)(c))
}

type SpaceServicesConfig1 struct{}

func (c *SpaceServicesConfig1) InitDefaults() *SpaceServicesConfig1 {
	return c
}

type MockQSpaceConfig1 struct {
	TokenAuth bool `json:"token_auth"`
}

func (c *MockQSpaceConfig1) InitDefaults() *MockQSpaceConfig1 {
	return c
}

type EthereumQSpaceConfig1 struct {
	NetworkId  uint64 `json:"network_id,omitempty"` // network ID
	URL        string `json:"url,omitempty"`        // URL to the blockchain
	WalletFile string `json:"wallet_file"`          // wallet file path
}

func (c *EthereumQSpaceConfig1) InitDefaults() *EthereumQSpaceConfig1 {
	c.NetworkId = 123
	c.WalletFile = "/predefined"
	return c
}

func (c *EthereumQSpaceConfig1) Validate() error {
	return nil
}

func (c *EthereumQSpaceConfig1) UnmarshalJSON(p []byte) error {
	c.InitDefaults()

	type alias EthereumQSpaceConfig1
	err := json.Unmarshal(p, (*alias)(c))
	if err == nil {
		//c.Bc.NetworkId = c.NetworkId
	}
	return err
}
