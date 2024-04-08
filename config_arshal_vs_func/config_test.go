package config_arshal_vs_func

import (
	"bytes"
	"fmt"
	"testing"

	jsonexp "github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/stretchr/testify/require"
)

var json1 = `{"qspaces":[{"id":"ispc123","names":[],"type":"Ethereum","ethereum":{"url":"slu://a/b"}}]}`

// TestUnmarshalJson unmarshal with encoding/json
func TestUnmarshalJson(t *testing.T) {
	c := &AppConfig1{}
	err := jsonexp.Unmarshal([]byte(json1), c)
	require.NoError(t, err)

	bb, err := jsonexp.Marshal(c)
	require.NoError(t, err)
	fmt.Println(string(bb))
	require.Equal(t, "ispc123", c.QSpaces[0].ID)
	require.Equal(t, "Ethereum", c.QSpaces[0].Type)
	require.NotNil(t, c.QSpaces[0].Ethereum)
	require.Equal(t, uint64(123), c.QSpaces[0].Ethereum.NetworkId)
	require.Equal(t, "slu://a/b", c.QSpaces[0].Ethereum.URL)
	require.Equal(t, "/predefined", c.QSpaces[0].Ethereum.WalletFile)
	require.NotNil(t, c.QSpaces[0].Mock)
}

// TestUnmarshalerV2 tests with UnmarshalerV2
func TestUnmarshalerV2(t *testing.T) {
	c := &AppConfig2{}
	err := jsonexp.Unmarshal([]byte(json1), c)
	require.NoError(t, err)

	bb, err := jsonexp.Marshal(c)
	require.NoError(t, err)
	fmt.Println(string(bb))
	require.Equal(t, "ispc123", c.QSpaces[0].ID)
	require.Equal(t, "Ethereum", c.QSpaces[0].Type)
	require.NotNil(t, c.QSpaces[0].Ethereum)
	require.Equal(t, uint64(123), c.QSpaces[0].Ethereum.NetworkId)
	require.Equal(t, "slu://a/b", c.QSpaces[0].Ethereum.URL)
	require.Equal(t, "/predefined", c.QSpaces[0].Ethereum.WalletFile)
	// 'mock' is now nil (different from encoding/json)
	require.Nil(t, c.QSpaces[0].Mock)
}

// TestUnmarshalFuncV2 tests with UnmarshalFuncV2
func TestUnmarshalFuncV2(t *testing.T) {
	c := &AppConfig3{}

	dec := jsontext.NewDecoder(bytes.NewReader([]byte(json1)))
	err := jsonexp.UnmarshalDecode(dec, c,
		jsonexp.WithUnmarshalers(
			jsonexp.UnmarshalFuncV2(func(dec *jsontext.Decoder, val any, opts jsonexp.Options) error {
				switch x := val.(type) {
				case *AppConfig3:
					fmt.Println("app config")
					x.InitDefaults()
				case *QSpacesConfig3:
					fmt.Println("qspace config")
					x.InitDefaults()
				case *EthereumQSpaceConfig3:
					fmt.Println("ethereum config")
					x.InitDefaults()
				}
				return jsonexp.SkipFunc
			})))
	require.NoError(t, err)

	bb, err := jsonexp.Marshal(c)
	require.NoError(t, err)
	fmt.Println(string(bb))
	require.Equal(t, "ispc123", c.QSpaces[0].ID)
	require.Equal(t, "Ethereum", c.QSpaces[0].Type)
	require.NotNil(t, c.QSpaces[0].Ethereum)
	require.Equal(t, uint64(123), c.QSpaces[0].Ethereum.NetworkId)
	require.Equal(t, "slu://a/b", c.QSpaces[0].Ethereum.URL)
	require.Equal(t, "/predefined", c.QSpaces[0].Ethereum.WalletFile)
	require.Nil(t, c.QSpaces[0].Mock)
}

// TestUnmarshalV2Skip tests that returning SkipFunc works with UnmarshalerV2
func TestUnmarshalV2Skip(t *testing.T) {
	c := &AppConfig4{}
	err := jsonexp.Unmarshal([]byte(json1), c)
	require.NoError(t, err)

	bb, err := jsonexp.Marshal(c)
	require.NoError(t, err)
	fmt.Println(string(bb))
	require.Equal(t, "ispc123", c.QSpaces[0].ID)
	require.Equal(t, "Ethereum", c.QSpaces[0].Type)
	require.NotNil(t, c.QSpaces[0].Ethereum)
	require.Equal(t, uint64(123), c.QSpaces[0].Ethereum.NetworkId)
	require.Equal(t, "slu://a/b", c.QSpaces[0].Ethereum.URL)
	require.Equal(t, "/predefined", c.QSpaces[0].Ethereum.WalletFile)
	require.Nil(t, c.QSpaces[0].Mock)
}

/*
=== RUN   TestUnmarshalJson
{"qspaces":[{"id":"ispc123","names":[],"type":"Ethereum","mock":{"token_auth":false},"ethereum":{"network_id":123,"url":"slu://a/b","wallet_file":"/predefined"},"services":{}}]}

=== RUN   TestUnmarshalerV2
{"qspaces":[{"id":"ispc123","names":[],"type":"Ethereum","ethereum":{"network_id":123,"url":"slu://a/b","wallet_file":"/predefined"},"services":{}}]}

=== RUN   TestUnmarshalFuncV2
{"qspaces":[{"id":"ispc123","names":[],"type":"Ethereum","mock":{"token_auth":false},"ethereum":{"network_id":123,"url":"slu://a/b","wallet_file":"/predefined"},"services":{}}]}

=== RUN   TestUnmarshalV2Skip
{"qspaces":[{"id":"ispc123","names":[],"type":"Ethereum","ethereum":{"network_id":123,"url":"slu://a/b","wallet_file":"/predefined"},"services":{}}]}


*/
