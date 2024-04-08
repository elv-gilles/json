package config_arshal_vs_func

import (
	"encoding/json"
	"fmt"

	jsonexp "github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type AppConfig4 struct {
	QSpacesConfig4
}

// InitDefaults initializes the default configuration.
func (c *AppConfig4) InitDefaults() *AppConfig4 {
	c.QSpacesConfig4.InitDefaults()
	return c
}

func (c *AppConfig4) UnmarshalJSON(bts []byte) error {
	c.InitDefaults()

	type alias AppConfig4
	return json.Unmarshal(bts, (*alias)(c))
}

func (c *AppConfig4) UnmarshalJSONV2(dec *jsontext.Decoder, opts jsonexp.Options) error {
	c.InitDefaults()
	return jsonexp.SkipFunc
}

type QSpacesConfig4 struct {
	QSpaces []QSpaceConfig4 `json:"qspaces"`
}

func (c *QSpacesConfig4) Validate() error {
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

func (c *QSpacesConfig4) InitDefaults() *QSpacesConfig4 {
	c.QSpaces = make([]QSpaceConfig4, 1)
	c.QSpaces[0] = QSpaceConfig4{
		ID:   "ispcMockContentSpace",
		Type: "Mock",
		Mock: (&MockQSpaceConfig4{}).InitDefaults(),
	}
	c.QSpaces[0].Services.InitDefaults()
	return c
}

type QSpaceConfig4 struct {
	ID       string                 `json:"id"`             // space ID
	Names    []string               `json:"names"`          // space names
	Type     string                 `json:"type"`           // Type is one of the supported content space types
	Mock     *MockQSpaceConfig4     `json:"mock,omitempty"` // Mock is the configuration for a mock content space
	Ethereum *EthereumQSpaceConfig4 `json:"ethereum,omitempty"`
	Services SpaceServicesConfig4   `json:"services"`
}

func (c *QSpaceConfig4) InitDefaults() *QSpaceConfig4 {
	c.Services.InitDefaults()
	return c
}

func (c *QSpaceConfig4) Validate() error {
	if c.Type == "Ethereum" {
		if c.Ethereum == nil {
			return fmt.Errorf("ethereum not defined")
		}
		return c.Ethereum.Validate()
	}
	return nil
}

func (c *QSpaceConfig4) UnmarshalJSON(bts []byte) error {
	c.InitDefaults()

	type alias QSpaceConfig4
	return json.Unmarshal(bts, (*alias)(c))
}

func (c *QSpaceConfig4) UnmarshalJSONV2(dec *jsontext.Decoder, opts jsonexp.Options) error {
	c.InitDefaults()
	return jsonexp.SkipFunc
}

type SpaceServicesConfig4 struct{}

func (c *SpaceServicesConfig4) InitDefaults() *SpaceServicesConfig4 {
	return c
}

type MockQSpaceConfig4 struct {
	TokenAuth bool `json:"token_auth"`
}

func (c *MockQSpaceConfig4) InitDefaults() *MockQSpaceConfig4 {
	return c
}

type EthereumQSpaceConfig4 struct {
	NetworkId  uint64 `json:"network_id,omitempty"` // network ID
	URL        string `json:"url,omitempty"`        // URL to the blockchain
	WalletFile string `json:"wallet_file"`          // wallet file path
}

func (c *EthereumQSpaceConfig4) InitDefaults() *EthereumQSpaceConfig4 {
	c.NetworkId = 123
	c.WalletFile = "/predefined"
	return c
}

func (c *EthereumQSpaceConfig4) Validate() error {
	return nil
}

func (c *EthereumQSpaceConfig4) UnmarshalJSON(p []byte) error {
	c.InitDefaults()

	type alias EthereumQSpaceConfig4
	err := json.Unmarshal(p, (*alias)(c))
	if err == nil {
		//c.Bc.NetworkId = c.NetworkId
	}
	return err
}

func (c *EthereumQSpaceConfig4) UnmarshalJSONV2(dec *jsontext.Decoder, opts jsonexp.Options) error {
	c.InitDefaults()
	return jsonexp.SkipFunc
}
