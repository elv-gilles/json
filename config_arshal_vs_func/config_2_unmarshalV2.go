package config_arshal_vs_func

import (
	"encoding/json"
	"fmt"

	jsonexp "github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type AppConfig2 struct {
	QSpacesConfig2
}

// InitDefaults initializes the default configuration.
func (c *AppConfig2) InitDefaults() *AppConfig2 {
	c.QSpacesConfig2.InitDefaults()
	return c
}

func (c *AppConfig2) UnmarshalJSON(bts []byte) error {
	c.InitDefaults()

	type alias AppConfig2
	return json.Unmarshal(bts, (*alias)(c))
}

func (c *AppConfig2) UnmarshalJSONV2(dec *jsontext.Decoder, opts jsonexp.Options) error {
	c.InitDefaults()

	type alias AppConfig2
	return jsonexp.UnmarshalDecode(dec, (*alias)(c), opts)
}

type QSpacesConfig2 struct {
	QSpaces []QSpaceConfig2 `json:"qspaces"`
}

func (c *QSpacesConfig2) Validate() error {
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

func (c *QSpacesConfig2) InitDefaults() *QSpacesConfig2 {
	c.QSpaces = make([]QSpaceConfig2, 1)
	c.QSpaces[0] = QSpaceConfig2{
		ID:   "ispcMockContentSpace",
		Type: "Mock",
		Mock: (&MockQSpaceConfig2{}).InitDefaults(),
	}
	c.QSpaces[0].Services.InitDefaults()
	return c
}

type QSpaceConfig2 struct {
	ID       string                 `json:"id"`             // space ID
	Names    []string               `json:"names"`          // space names
	Type     string                 `json:"type"`           // Type is one of the supported content space types
	Mock     *MockQSpaceConfig2     `json:"mock,omitempty"` // Mock is the configuration for a mock content space
	Ethereum *EthereumQSpaceConfig2 `json:"ethereum,omitempty"`
	Services SpaceServicesConfig2   `json:"services"`
}

func (c *QSpaceConfig2) InitDefaults() *QSpaceConfig2 {
	c.Services.InitDefaults()
	return c
}

func (c *QSpaceConfig2) Validate() error {
	if c.Type == "Ethereum" {
		if c.Ethereum == nil {
			return fmt.Errorf("ethereum not defined")
		}
		return c.Ethereum.Validate()
	}
	return nil
}

func (c *QSpaceConfig2) UnmarshalJSON(bts []byte) error {
	c.InitDefaults()

	type alias QSpaceConfig2
	return json.Unmarshal(bts, (*alias)(c))
}

func (c *QSpaceConfig2) UnmarshalJSONV2(dec *jsontext.Decoder, opts jsonexp.Options) error {
	c.InitDefaults()

	type alias QSpaceConfig2
	return jsonexp.UnmarshalDecode(dec, (*alias)(c), opts)
}

type SpaceServicesConfig2 struct{}

func (c *SpaceServicesConfig2) InitDefaults() *SpaceServicesConfig2 {
	return c
}

type MockQSpaceConfig2 struct {
	TokenAuth bool `json:"token_auth"`
}

func (c *MockQSpaceConfig2) InitDefaults() *MockQSpaceConfig2 {
	return c
}

type EthereumQSpaceConfig2 struct {
	NetworkId  uint64 `json:"network_id,omitempty"` // network ID
	URL        string `json:"url,omitempty"`        // URL to the blockchain
	WalletFile string `json:"wallet_file"`          // wallet file path
}

func (c *EthereumQSpaceConfig2) InitDefaults() *EthereumQSpaceConfig2 {
	c.NetworkId = 123
	c.WalletFile = "/predefined"
	return c
}

func (c *EthereumQSpaceConfig2) Validate() error {
	return nil
}

func (c *EthereumQSpaceConfig2) UnmarshalJSON(p []byte) error {
	c.InitDefaults()

	type alias EthereumQSpaceConfig2
	err := json.Unmarshal(p, (*alias)(c))
	if err == nil {
		//c.Bc.NetworkId = c.NetworkId
	}
	return err
}

func (c *EthereumQSpaceConfig2) UnmarshalJSONV2(dec *jsontext.Decoder, opts jsonexp.Options) error {
	c.InitDefaults()

	type alias EthereumQSpaceConfig2
	err := jsonexp.UnmarshalDecode(dec, (*alias)(c), opts)
	if err == nil {
		//c.Bc.NetworkId = c.NetworkId
	}
	return err
}
