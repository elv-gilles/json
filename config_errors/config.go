package config

import (
	"encoding/json"
	"fmt"

	jsonexp "github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type AppConfig struct {
	// QSpacesConfig is the configuration of content spaces
	QSpacesConfig
}

// InitDefaults initializes the default configuration.
func (c *AppConfig) InitDefaults() *AppConfig {
	c.QSpacesConfig.InitDefaults()
	return c
}

func (c *AppConfig) UnmarshalJSON(bts []byte) error {
	c.InitDefaults()

	type alias AppConfig
	return json.Unmarshal(bts, (*alias)(c))
}

func (c *AppConfig) UnmarshalJSONV2(dec *jsontext.Decoder, opts jsonexp.Options) error {
	c.InitDefaults()
	type alias AppConfig
	return jsonexp.UnmarshalDecode(dec, (*alias)(c), opts)
}

// QSpacesConfig is the configuration for the supported content spaces of this fabric node.
type QSpacesConfig struct {
	QSpaces []QSpaceConfig `json:"qspaces"`
}

func (c *QSpacesConfig) Validate() error {
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

func (c *QSpacesConfig) InitDefaults() *QSpacesConfig {
	c.QSpaces = make([]QSpaceConfig, 1)
	c.QSpaces[0] = QSpaceConfig{
		ID:   "ispcMockContentSpace",
		Type: "Mock",
		Mock: (&MockQSpaceConfig{}).InitDefaults(),
	}
	c.QSpaces[0].Services.InitDefaults()
	return c
}

// QSpaceConfig is the configuration of a content space - see qspaceConfigFields.
type QSpaceConfig struct {
	// ID is the identifier of the content space
	ID       string                `json:"id"`             // space ID
	Names    []string              `json:"names"`          // space names
	Type     string                `json:"type"`           // Type is one of the supported content space types
	Mock     *MockQSpaceConfig     `json:"mock,omitempty"` // Mock is the configuration for a mock content space
	Ethereum *EthereumQSpaceConfig `json:"ethereum,omitempty"`
	Services SpaceServicesConfig   `json:"services"`
}

func (c *QSpaceConfig) InitDefaults() *QSpaceConfig {
	c.Services.InitDefaults()
	return c
}

func (c *QSpaceConfig) Validate() error {
	if c.Type == "Ethereum" {
		if c.Ethereum == nil {
			return fmt.Errorf("ethereum not defined")
		}
		return c.Ethereum.Validate()
	}
	return nil
}

func (c *QSpaceConfig) UnmarshalJSON(bts []byte) error {
	c.InitDefaults()

	type alias QSpaceConfig
	return json.Unmarshal(bts, (*alias)(c))
}

func (c *QSpaceConfig) UnmarshalJSONV2(dec *jsontext.Decoder, opts jsonexp.Options) error {
	c.InitDefaults()
	type alias QSpaceConfig
	return jsonexp.UnmarshalDecode(dec, (*alias)(c), opts)
}

type SpaceServicesConfig struct{}

func (c *SpaceServicesConfig) InitDefaults() *SpaceServicesConfig {
	return c
}

type MockQSpaceConfig struct {
	TokenAuth bool `json:"token_auth"`
}

func (c *MockQSpaceConfig) InitDefaults() *MockQSpaceConfig {
	return c
}

type EthereumQSpaceConfig struct {
	NetworkId  uint64 `json:"network_id,omitempty"` // network ID
	URL        string `json:"url,omitempty"`        // URL to the blockchain
	WalletFile string `json:"wallet_file"`          // wallet file path
}

func (c *EthereumQSpaceConfig) InitDefaults() *EthereumQSpaceConfig {
	return c
}

func (c *EthereumQSpaceConfig) Validate() error {
	return nil
}

func (c *EthereumQSpaceConfig) UnmarshalJSON(p []byte) error {
	c.InitDefaults()

	type ethCfg EthereumQSpaceConfig
	var o = (*ethCfg)(c)

	err := json.Unmarshal(p, o)
	if err == nil {
		//c.Bc.NetworkId = c.NetworkId
	}
	return err
}

func (c *EthereumQSpaceConfig) UnmarshalJSONV2(dec *jsontext.Decoder, opts jsonexp.Options) error {
	c.InitDefaults()
	type alias EthereumQSpaceConfig

	err := jsonexp.UnmarshalDecode(dec, (*alias)(c), opts)
	if err == nil {
		//c.Bc.NetworkId = c.NetworkId
	}
	return err
}
