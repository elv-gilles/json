package config

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	jsonexp "github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	emptyJsonConfig  = `{"qspaces":[{"id":"ispcMockContentSpace","names":null,"type":"Mock","mock":{"token_auth":false},"services":{}}]}`
	simpleJsonConfig = `{
 "qspaces": [
  {
   "id": "eth_123",
   "names": [
    "x"
   ],
   "type": "Ethereum",
   "ethereum": {
    "network_id": 123,
    "url": "sss://xyz",
    "wallet_file": "/a/b/c"
   },
   "services": {}
  }
 ]
}`
)

func simpleConfig() *AppConfig {
	c := (&AppConfig{}).InitDefaults()
	c.QSpaces[0] = QSpaceConfig{
		ID:    "eth_123",
		Names: []string{"x"},
		Type:  "Ethereum",
		Mock:  (&MockQSpaceConfig{}).InitDefaults(),
		Ethereum: &EthereumQSpaceConfig{
			NetworkId:  123,
			URL:        "sss://xyz",
			WalletFile: "/a/b/c",
		},
		Services: SpaceServicesConfig{},
	}
	return c
}

func TestUnmarshalJSON(t *testing.T) {
	c := &AppConfig{}
	err := json.Unmarshal([]byte(simpleJsonConfig), c)
	require.NoError(t, err)

	require.EqualValues(t, simpleConfig(), c)
}

var jsonTemplateNetworkId = `{ "qspaces": [{
   		"id": "eth_123",
   		"names": [    "x"   ],
   		"type": "Ethereum",
   		"ethereum": {
   		 "network_id": $$,
   		 "url": "sss://xyz",
   		 "wallet_file": "/a/b/c"
   	}}]}`

func TestUnmarshalJSONError(t *testing.T) {
	type testCase struct {
		networkId string
		wantErr   string
	}

	for i, tc := range []*testCase{
		{networkId: "123x", wantErr: "invalid character 'x' after object key:value pair"},
		{networkId: "-123", wantErr: "json: cannot unmarshal number -123 into Go struct field alias.qspaces of type uint64"},
		{networkId: "\"123:456\"", wantErr: "json: cannot unmarshal string into Go struct field alias.qspaces of type uint64"},
	} {
		c := &AppConfig{}
		js := strings.Replace(jsonTemplateNetworkId, "$$", tc.networkId, 1)
		err := json.Unmarshal([]byte(js), c)
		require.Error(t, err, "test-case %d", i)
		require.Equal(t, tc.wantErr, err.Error(), "test-case %d", i)
	}
}

func printError(err error, lvl int) {
	prefix := ""
	for len(prefix) < lvl {
		prefix = prefix + "  "
	}
	switch errx := err.(type) {
	case *jsonexp.SemanticError:
		fmt.Println(prefix+"semantic",
			"go-type", errx.GoType,
			"offset", errx.ByteOffset,
			"json_kind", errx.JSONKind,
			"json_pointer", errx.JSONPointer)
		printError(errx.Err, lvl+1)
	case *jsontext.SyntacticError:
		fmt.Println(prefix+"syntax", "offset", errx.ByteOffset)
	default:
		return
	}
}

func TestUnmarshalJSONV2Error(t *testing.T) {
	type testCase struct {
		networkId string
		wantErr   string
	}

	unrender := func(s string) string {
		return strings.ReplaceAll(s, "unable to", "cannot")
	}
	assertError := func(tc *testCase, err error, i int) {
		assert.Equal(t, unrender(tc.wantErr), unrender(err.Error()), "test-case %d", i)
	}

	for i, tc := range []*testCase{
		{
			networkId: "123x",
			wantErr:   `json: cannot unmarshal Go value of type config.EthereumQSpaceConfig within JSON value at "/qspaces/0/ethereum/network_id": jsontext: missing character ',' after object or array value at byte offset 133`,
		},
		{
			networkId: "-123",
			wantErr:   `json: cannot unmarshal JSON number into Go value of type uint64 within JSON value at "/qspaces/0/ethereum/network_id" at byte offset 134: cannot parse "-123" as unsigned integer: invalid syntax`,
		},
		{
			networkId: "\"123:456\"",
			wantErr:   `json: cannot unmarshal JSON string into Go value of type uint64 within JSON value at "/qspaces/0/ethereum/network_id" at byte offset 139: invalid value: "123:456"`,
		},
	} {
		c := &AppConfig{}
		js := strings.Replace(jsonTemplateNetworkId, "$$", tc.networkId, 1)
		err := jsonexp.Unmarshal([]byte(js), c)
		assert.Error(t, err, "test-case %d", i)
		//printError(err, 0)

		assertError(tc, err, i)
	}
}
