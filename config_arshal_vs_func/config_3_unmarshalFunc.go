package config_arshal_vs_func

import (
	"fmt"
)

type AppConfig3 struct {
	QSpacesConfig3
}

// InitDefaults initializes the default configuration.
func (c *AppConfig3) InitDefaults() *AppConfig3 {
	c.QSpacesConfig3.InitDefaults()
	return c
}

//func (c *AppConfig3) UnmarshalJSON(bts []byte) error {
//	c.InitDefaults()
//
//	type alias AppConfig3
//	return json.Unmarshal(bts, (*alias)(c))
//}

type QSpacesConfig3 struct {
	QSpaces []QSpaceConfig3 `json:"qspaces"`
}

func (c *QSpacesConfig3) Validate() error {
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

func (c *QSpacesConfig3) InitDefaults() *QSpacesConfig3 {
	c.QSpaces = make([]QSpaceConfig3, 1)
	c.QSpaces[0] = QSpaceConfig3{
		ID:   "ispcMockContentSpace",
		Type: "Mock",
		Mock: (&MockQSpaceConfig3{}).InitDefaults(),
	}
	c.QSpaces[0].Services.InitDefaults()
	return c
}

type QSpaceConfig3 struct {
	ID       string                 `json:"id"`             // space ID
	Names    []string               `json:"names"`          // space names
	Type     string                 `json:"type"`           // Type is one of the supported content space types
	Mock     *MockQSpaceConfig3     `json:"mock,omitempty"` // Mock is the configuration for a mock content space
	Ethereum *EthereumQSpaceConfig3 `json:"ethereum,omitempty"`
	Services SpaceServicesConfig3   `json:"services"`
}

func (c *QSpaceConfig3) InitDefaults() *QSpaceConfig3 {
	c.Services.InitDefaults()
	return c
}

func (c *QSpaceConfig3) Validate() error {
	if c.Type == "Ethereum" {
		if c.Ethereum == nil {
			return fmt.Errorf("ethereum not defined")
		}
		return c.Ethereum.Validate()
	}
	return nil
}

//func (c *QSpaceConfig3) UnmarshalJSON(bts []byte) error {
//	c.InitDefaults()
//
//	type alias QSpaceConfig3
//	return json.Unmarshal(bts, (*alias)(c))
//}

type SpaceServicesConfig3 struct{}

func (c *SpaceServicesConfig3) InitDefaults() *SpaceServicesConfig3 {
	return c
}

type MockQSpaceConfig3 struct {
	TokenAuth bool `json:"token_auth"`
}

func (c *MockQSpaceConfig3) InitDefaults() *MockQSpaceConfig3 {
	return c
}

type EthereumQSpaceConfig3 struct {
	NetworkId  uint64 `json:"network_id,omitempty"` // network ID
	URL        string `json:"url,omitempty"`        // URL to the blockchain
	WalletFile string `json:"wallet_file"`          // wallet file path
}

func (c *EthereumQSpaceConfig3) InitDefaults() *EthereumQSpaceConfig3 {
	c.NetworkId = 123
	c.WalletFile = "/predefined"
	return c
}

func (c *EthereumQSpaceConfig3) Validate() error {
	return nil
}

//func (c *EthereumQSpaceConfig3) UnmarshalJSON(p []byte) error {
//	c.InitDefaults()
//
//	type alias EthereumQSpaceConfig3
//	err := json.Unmarshal(p, (*alias)(c))
//	if err == nil {
//		//c.Bc.NetworkId = c.NetworkId
//	}
//	return err
//}
