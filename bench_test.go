// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package json

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"
	"testing/iotest"
	"time"

	jsonv1 "encoding/json"

	"github.com/go-json-experiment/json/internal/jsontest"
	"github.com/go-json-experiment/json/jsontext"
)

var benchV1 = os.Getenv("BENCHMARK_V1") != ""

// bytesBuffer is identical to bytes.Buffer,
// but a different type to avoid any optimizations for bytes.Buffer.
type bytesBuffer struct{ *bytes.Buffer }

var arshalTestdata = []struct {
	name   string
	raw    []byte
	val    any
	new    func() any
	skipV1 bool
}{{
	name: "Bool",
	raw:  []byte("true"),
	val:  addr(true),
	new:  func() any { return new(bool) },
}, {
	name: "String",
	raw:  []byte(`"hello, world!"`),
	val:  addr("hello, world!"),
	new:  func() any { return new(string) },
}, {
	name: "Int",
	raw:  []byte("-1234"),
	val:  addr(int64(-1234)),
	new:  func() any { return new(int64) },
}, {
	name: "Uint",
	raw:  []byte("1234"),
	val:  addr(uint64(1234)),
	new:  func() any { return new(uint64) },
}, {
	name: "Float",
	raw:  []byte("12.34"),
	val:  addr(float64(12.34)),
	new:  func() any { return new(float64) },
}, {
	name: "Map/ManyEmpty",
	raw:  []byte(`[{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{}]`),
	val: addr(func() (out []map[string]string) {
		for range 100 {
			out = append(out, map[string]string{})
		}
		return out
	}()),
	new: func() any { return new([]map[string]string) },
}, {
	name: "Map/OneLarge",
	raw:  []byte(`{"A":"A","B":"B","C":"C","D":"D","E":"E","F":"F","G":"G","H":"H","I":"I","J":"J","K":"K","L":"L","M":"M","N":"N","O":"O","P":"P","Q":"Q","R":"R","S":"S","T":"T","U":"U","V":"V","W":"W","X":"X","Y":"Y","Z":"Z"}`),
	val:  addr(map[string]string{"A": "A", "B": "B", "C": "C", "D": "D", "E": "E", "F": "F", "G": "G", "H": "H", "I": "I", "J": "J", "K": "K", "L": "L", "M": "M", "N": "N", "O": "O", "P": "P", "Q": "Q", "R": "R", "S": "S", "T": "T", "U": "U", "V": "V", "W": "W", "X": "X", "Y": "Y", "Z": "Z"}),
	new:  func() any { return new(map[string]string) },
}, {
	name: "Map/ManySmall",
	raw:  []byte(`{"A":{"K":"V"},"B":{"K":"V"},"C":{"K":"V"},"D":{"K":"V"},"E":{"K":"V"},"F":{"K":"V"},"G":{"K":"V"},"H":{"K":"V"},"I":{"K":"V"},"J":{"K":"V"},"K":{"K":"V"},"L":{"K":"V"},"M":{"K":"V"},"N":{"K":"V"},"O":{"K":"V"},"P":{"K":"V"},"Q":{"K":"V"},"R":{"K":"V"},"S":{"K":"V"},"T":{"K":"V"},"U":{"K":"V"},"V":{"K":"V"},"W":{"K":"V"},"X":{"K":"V"},"Y":{"K":"V"},"Z":{"K":"V"}}`),
	val:  addr(map[string]map[string]string{"A": {"K": "V"}, "B": {"K": "V"}, "C": {"K": "V"}, "D": {"K": "V"}, "E": {"K": "V"}, "F": {"K": "V"}, "G": {"K": "V"}, "H": {"K": "V"}, "I": {"K": "V"}, "J": {"K": "V"}, "K": {"K": "V"}, "L": {"K": "V"}, "M": {"K": "V"}, "N": {"K": "V"}, "O": {"K": "V"}, "P": {"K": "V"}, "Q": {"K": "V"}, "R": {"K": "V"}, "S": {"K": "V"}, "T": {"K": "V"}, "U": {"K": "V"}, "V": {"K": "V"}, "W": {"K": "V"}, "X": {"K": "V"}, "Y": {"K": "V"}, "Z": {"K": "V"}}),
	new:  func() any { return new(map[string]map[string]string) },
}, {
	name: "Struct/ManyEmpty",
	raw:  []byte(`[{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{}]`),
	val:  addr(make([]struct{}, 100)),
	new: func() any {
		return new([]struct{})
	},
}, {
	name: "Struct/OneLarge",
	raw:  []byte(`{"A":"A","B":"B","C":"C","D":"D","E":"E","F":"F","G":"G","H":"H","I":"I","J":"J","K":"K","L":"L","M":"M","N":"N","O":"O","P":"P","Q":"Q","R":"R","S":"S","T":"T","U":"U","V":"V","W":"W","X":"X","Y":"Y","Z":"Z"}`),
	val:  addr(struct{ A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R, S, T, U, V, W, X, Y, Z string }{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}),
	new: func() any {
		return new(struct{ A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R, S, T, U, V, W, X, Y, Z string })
	},
}, {
	name: "Struct/ManySmall",
	raw:  []byte(`{"A":{"K":"V"},"B":{"K":"V"},"C":{"K":"V"},"D":{"K":"V"},"E":{"K":"V"},"F":{"K":"V"},"G":{"K":"V"},"H":{"K":"V"},"I":{"K":"V"},"J":{"K":"V"},"K":{"K":"V"},"L":{"K":"V"},"M":{"K":"V"},"N":{"K":"V"},"O":{"K":"V"},"P":{"K":"V"},"Q":{"K":"V"},"R":{"K":"V"},"S":{"K":"V"},"T":{"K":"V"},"U":{"K":"V"},"V":{"K":"V"},"W":{"K":"V"},"X":{"K":"V"},"Y":{"K":"V"},"Z":{"K":"V"}}`),
	val: func() any {
		V := struct{ K string }{"V"}
		return addr(struct{ A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R, S, T, U, V, W, X, Y, Z struct{ K string } }{
			V, V, V, V, V, V, V, V, V, V, V, V, V, V, V, V, V, V, V, V, V, V, V, V, V, V,
		})
	}(),
	new: func() any {
		return new(struct{ A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R, S, T, U, V, W, X, Y, Z struct{ K string } })
	},
}, {
	name: "Slice/ManyEmpty",
	raw:  []byte(`[[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[]]`),
	val: addr(func() (out [][]string) {
		for range 100 {
			out = append(out, []string{})
		}
		return out
	}()),
	new: func() any { return new([][]string) },
}, {
	name: "Slice/OneLarge",
	raw:  []byte(`["A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z"]`),
	val:  addr([]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}),
	new:  func() any { return new([]string) },
}, {
	name: "Slice/ManySmall",
	raw:  []byte(`[["A"],["B"],["C"],["D"],["E"],["F"],["G"],["H"],["I"],["J"],["K"],["L"],["M"],["N"],["O"],["P"],["Q"],["R"],["S"],["T"],["U"],["V"],["W"],["X"],["Y"],["Z"]]`),
	val:  addr([][]string{{"A"}, {"B"}, {"C"}, {"D"}, {"E"}, {"F"}, {"G"}, {"H"}, {"I"}, {"J"}, {"K"}, {"L"}, {"M"}, {"N"}, {"O"}, {"P"}, {"Q"}, {"R"}, {"S"}, {"T"}, {"U"}, {"V"}, {"W"}, {"X"}, {"Y"}, {"Z"}}),
	new:  func() any { return new([][]string) },
}, {
	name: "Array/OneLarge",
	raw:  []byte(`["A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z"]`),
	val:  addr([26]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}),
	new:  func() any { return new([26]string) },
}, {
	name: "Array/ManySmall",
	raw:  []byte(`[["A"],["B"],["C"],["D"],["E"],["F"],["G"],["H"],["I"],["J"],["K"],["L"],["M"],["N"],["O"],["P"],["Q"],["R"],["S"],["T"],["U"],["V"],["W"],["X"],["Y"],["Z"]]`),
	val:  addr([26][1]string{{"A"}, {"B"}, {"C"}, {"D"}, {"E"}, {"F"}, {"G"}, {"H"}, {"I"}, {"J"}, {"K"}, {"L"}, {"M"}, {"N"}, {"O"}, {"P"}, {"Q"}, {"R"}, {"S"}, {"T"}, {"U"}, {"V"}, {"W"}, {"X"}, {"Y"}, {"Z"}}),
	new:  func() any { return new([26][1]string) },
}, {
	name: "Bytes/Slice",
	raw:  []byte(`"47DEQpj8HBSa+/TImW+5JCeuQeRkm5NMpJWZG3hSuFU="`),
	val:  addr([]byte{0xe3, 0xb0, 0xc4, 0x42, 0x98, 0xfc, 0x1c, 0x14, 0x9a, 0xfb, 0xf4, 0xc8, 0x99, 0x6f, 0xb9, 0x24, 0x27, 0xae, 0x41, 0xe4, 0x64, 0x9b, 0x93, 0x4c, 0xa4, 0x95, 0x99, 0x1b, 0x78, 0x52, 0xb8, 0x55}),
	new:  func() any { return new([]byte) },
}, {
	name:   "Bytes/Array",
	raw:    []byte(`"47DEQpj8HBSa+/TImW+5JCeuQeRkm5NMpJWZG3hSuFU="`),
	val:    addr([32]byte{0xe3, 0xb0, 0xc4, 0x42, 0x98, 0xfc, 0x1c, 0x14, 0x9a, 0xfb, 0xf4, 0xc8, 0x99, 0x6f, 0xb9, 0x24, 0x27, 0xae, 0x41, 0xe4, 0x64, 0x9b, 0x93, 0x4c, 0xa4, 0x95, 0x99, 0x1b, 0x78, 0x52, 0xb8, 0x55}),
	new:    func() any { return new([32]byte) },
	skipV1: true,
}, {
	name: "Pointer",
	raw:  []byte("true"),
	val:  addr(addr(addr(addr(addr(addr(addr(addr(addr(addr(addr(true))))))))))),
	new:  func() any { return new(**********bool) },
}, {
	name: "TextArshal",
	raw:  []byte(`"method"`),
	val:  new(textArshaler),
	new:  func() any { return new(textArshaler) },
}, {
	name: "JSONArshalV1",
	raw:  []byte(`"method"`),
	val:  new(jsonArshalerV1),
	new:  func() any { return new(jsonArshalerV1) },
}, {
	name:   "JSONArshalV2",
	raw:    []byte(`"method"`),
	val:    new(jsonArshalerV2),
	new:    func() any { return new(jsonArshalerV2) },
	skipV1: true,
}, {
	name:   "Duration",
	raw:    []byte(`"1h1m1s"`),
	val:    addr(time.Hour + time.Minute + time.Second),
	new:    func() any { return new(time.Duration) },
	skipV1: true,
}, {
	name: "Time",
	raw:  []byte(`"2006-01-02T22:04:05Z"`),
	val:  addr(time.Unix(1136239445, 0).UTC()),
	new:  func() any { return new(time.Time) },
}}

type textArshaler struct{ _ [4]int }

func (textArshaler) MarshalText() ([]byte, error) {
	return []byte("method"), nil
}
func (*textArshaler) UnmarshalText(b []byte) error {
	if string(b) != "method" {
		return fmt.Errorf("UnmarshalText: got %q, want %q", b, "method")
	}
	return nil
}

type jsonArshalerV1 struct{ _ [4]int }

func (jsonArshalerV1) MarshalJSON() ([]byte, error) {
	return []byte(`"method"`), nil
}
func (*jsonArshalerV1) UnmarshalJSON(b []byte) error {
	if string(b) != `"method"` {
		return fmt.Errorf("UnmarshalJSON: got %q, want %q", b, `"method"`)
	}
	return nil
}

type jsonArshalerV2 struct{ _ [4]int }

func (jsonArshalerV2) MarshalJSONV2(enc *jsontext.Encoder, opts Options) error {
	return enc.WriteToken(jsontext.String("method"))
}
func (*jsonArshalerV2) UnmarshalJSONV2(dec *jsontext.Decoder, opts Options) error {
	b, err := dec.ReadValue()
	if string(b) != `"method"` {
		return fmt.Errorf("UnmarshalJSONV2: got %q, want %q", b, `"method"`)
	}
	return err
}

func TestBenchmarkUnmarshal(t *testing.T) { runUnmarshal(t) }
func BenchmarkUnmarshal(b *testing.B)     { runUnmarshal(b) }

func runUnmarshal(tb testing.TB) {
	for _, tt := range arshalTestdata {
		// Setup the unmarshal operation.
		var val any
		run := func(tb testing.TB) {
			val = tt.new()
			if err := Unmarshal(tt.raw, val); err != nil {
				tb.Fatalf("Unmarshal error: %v", err)
			}
		}
		if benchV1 {
			run = func(tb testing.TB) {
				if tt.skipV1 {
					tb.Skip("not supported in v1")
				}
				val = tt.new()
				if err := jsonv1.Unmarshal(tt.raw, val); err != nil {
					tb.Fatalf("Unmarshal error: %v", err)
				}
			}
		}

		// Verify the results.
		if _, ok := tb.(*testing.T); ok {
			run0 := run
			run = func(tb testing.TB) {
				run0(tb)
				if !reflect.DeepEqual(val, tt.val) {
					tb.Fatalf("Unmarshal output mismatch:\ngot  %v\nwant %v", val, tt.val)
				}
			}
		}

		runTestOrBench(tb, tt.name, len64(tt.raw), run)
	}
}

func TestBenchmarkMarshal(t *testing.T) { runMarshal(t) }
func BenchmarkMarshal(b *testing.B)     { runMarshal(b) }

func runMarshal(tb testing.TB) {
	for _, tt := range arshalTestdata {
		// Setup the marshal operation.
		var raw []byte
		run := func(tb testing.TB) {
			var err error
			raw, err = Marshal(tt.val)
			if err != nil {
				tb.Fatalf("Marshal error: %v", err)
			}
		}
		if benchV1 {
			run = func(tb testing.TB) {
				if tt.skipV1 {
					tb.Skip("not supported in v1")
				}
				var err error
				raw, err = jsonv1.Marshal(tt.val)
				if err != nil {
					tb.Fatalf("Marshal error: %v", err)
				}
			}
		}

		// Verify the results.
		if _, ok := tb.(*testing.T); ok {
			run0 := run
			run = func(tb testing.TB) {
				run0(tb)
				if !bytes.Equal(raw, tt.raw) {
					// Map marshaling in v2 is non-deterministic.
					byteHistogram := func(b []byte) (h [256]int) {
						for _, c := range b {
							h[c]++
						}
						return h
					}
					if !(strings.HasPrefix(tt.name, "Map/") && byteHistogram(raw) == byteHistogram(tt.raw)) {
						tb.Fatalf("Marshal output mismatch:\ngot  %s\nwant %s", raw, tt.raw)
					}
				}
			}
		}

		runTestOrBench(tb, tt.name, len64(tt.raw), run)
	}
}

func TestBenchmarkTestdata(t *testing.T) { runAllTestdata(t) }
func BenchmarkTestdata(b *testing.B)     { runAllTestdata(b) }

func runAllTestdata(tb testing.TB) {
	for _, td := range jsontest.Data {
		for _, arshalName := range []string{"Marshal", "Unmarshal"} {
			for _, typeName := range []string{"Concrete", "Interface"} {
				newValue := func() any { return new(any) }
				if typeName == "Concrete" {
					if td.New == nil {
						continue
					}
					newValue = td.New
				}
				value := mustUnmarshalValue(tb, td.Data(), newValue)
				name := path.Join(td.Name, arshalName, typeName)
				runTestOrBench(tb, name, int64(td.Size), func(tb testing.TB) {
					runArshal(tb, arshalName, newValue, td.Data(), value)
				})
			}
		}

		tokens := mustDecodeTokens(tb, td.Data())
		buffer := make([]byte, 0, 2*len(td.Data()))
		for _, codeName := range []string{"Encode", "Decode"} {
			for _, typeName := range []string{"Token", "Value"} {
				for _, modeName := range []string{"Streaming", "Buffered"} {
					name := path.Join(td.Name, codeName, typeName, modeName)
					runTestOrBench(tb, name, int64(td.Size), func(tb testing.TB) {
						runCode(tb, codeName, typeName, modeName, buffer, td.Data(), tokens)
					})
				}
			}
		}
	}
}

func mustUnmarshalValue(t testing.TB, data []byte, newValue func() any) (value any) {
	value = newValue()
	if err := Unmarshal(data, value); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	return value
}

func runArshal(t testing.TB, arshalName string, newValue func() any, data []byte, value any) {
	if benchV1 {
		switch arshalName {
		case "Marshal":
			if _, err := jsonv1.Marshal(value); err != nil {
				t.Fatalf("Marshal error: %v", err)
			}
		case "Unmarshal":
			if err := jsonv1.Unmarshal(data, newValue()); err != nil {
				t.Fatalf("Unmarshal error: %v", err)
			}
		}
		return
	}

	switch arshalName {
	case "Marshal":
		if _, err := Marshal(value); err != nil {
			t.Fatalf("Marshal error: %v", err)
		}
	case "Unmarshal":
		if err := Unmarshal(data, newValue()); err != nil {
			t.Fatalf("Unmarshal error: %v", err)
		}
	}
}

func mustDecodeTokens(t testing.TB, data []byte) []jsontext.Token {
	var tokens []jsontext.Token
	dec := jsontext.NewDecoder(bytes.NewReader(data))
	for {
		tok, err := dec.ReadToken()
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatalf("Decoder.ReadToken error: %v", err)
		}

		// Prefer exact representation for JSON strings and numbers
		// since this more closely matches common use cases.
		switch tok.Kind() {
		case '"':
			tokens = append(tokens, jsontext.String(tok.String()))
		case '0':
			tokens = append(tokens, jsontext.Float(tok.Float()))
		default:
			tokens = append(tokens, tok.Clone())
		}
	}
	return tokens
}

func runCode(t testing.TB, codeName, typeName, modeName string, buffer, data []byte, tokens []jsontext.Token) {
	switch codeName {
	case "Encode":
		runEncode(t, typeName, modeName, buffer, data, tokens)
	case "Decode":
		runDecode(t, typeName, modeName, buffer, data, tokens)
	}
}

func runEncode(t testing.TB, typeName, modeName string, buffer, data []byte, tokens []jsontext.Token) {
	if benchV1 {
		if modeName == "Buffered" {
			t.Skip("no support for direct buffered input in v1")
		}
		enc := jsonv1.NewEncoder(bytes.NewBuffer(buffer[:0]))
		switch typeName {
		case "Token":
			t.Skip("no support for encoding tokens in v1; see https://go.dev/issue/40127")
		case "Value":
			val := jsonv1.RawMessage(data)
			if err := enc.Encode(val); err != nil {
				t.Fatalf("Decoder.Encode error: %v", err)
			}
		}
		return
	}

	var enc *jsontext.Encoder
	switch modeName {
	case "Streaming":
		enc = jsontext.NewEncoder(bytesBuffer{bytes.NewBuffer(buffer[:0])})
	case "Buffered":
		enc = jsontext.NewEncoder(bytes.NewBuffer(buffer[:0]))
	}
	switch typeName {
	case "Token":
		for _, tok := range tokens {
			if err := enc.WriteToken(tok); err != nil {
				t.Fatalf("Encoder.WriteToken error: %v", err)
			}
		}
	case "Value":
		if err := enc.WriteValue(data); err != nil {
			t.Fatalf("Encoder.WriteValue error: %v", err)
		}
	}
}

func runDecode(t testing.TB, typeName, modeName string, buffer, data []byte, tokens []jsontext.Token) {
	if benchV1 {
		if modeName == "Buffered" {
			t.Skip("no support for direct buffered output in v1")
		}
		dec := jsonv1.NewDecoder(bytes.NewReader(data))
		switch typeName {
		case "Token":
			for {
				if _, err := dec.Token(); err != nil {
					if err == io.EOF {
						break
					}
					t.Fatalf("Decoder.Token error: %v", err)
				}
			}
		case "Value":
			var val jsonv1.RawMessage
			if err := dec.Decode(&val); err != nil {
				t.Fatalf("Decoder.Decode error: %v", err)
			}
		}
		return
	}

	var dec *jsontext.Decoder
	switch modeName {
	case "Streaming":
		dec = jsontext.NewDecoder(bytesBuffer{bytes.NewBuffer(data)})
	case "Buffered":
		dec = jsontext.NewDecoder(bytes.NewBuffer(data))
	}
	switch typeName {
	case "Token":
		for {
			if _, err := dec.ReadToken(); err != nil {
				if err == io.EOF {
					break
				}
				t.Fatalf("Decoder.ReadToken error: %v", err)
			}
		}
	case "Value":
		if _, err := dec.ReadValue(); err != nil {
			t.Fatalf("Decoder.ReadValue error: %v", err)
		}
	}
}

var ws = strings.Repeat(" ", 4<<10)
var slowStreamingDecodeTestdata = []struct {
	name string
	data []byte
}{
	{"LargeString", []byte(`"` + strings.Repeat(" ", 4<<10) + `"`)},
	{"LargeNumber", []byte("0." + strings.Repeat("0", 4<<10))},
	{"LargeWhitespace/Null", []byte(ws + "null" + ws)},
	{"LargeWhitespace/Object", []byte(ws + "{" + ws + `"name1"` + ws + ":" + ws + `"value"` + ws + "," + ws + `"name2"` + ws + ":" + ws + `"value"` + ws + "}" + ws)},
	{"LargeWhitespace/Array", []byte(ws + "[" + ws + `"value"` + ws + "," + ws + `"value"` + ws + "]" + ws)},
}

func TestBenchmarkSlowStreamingDecode(t *testing.T) { runAllSlowStreamingDecode(t) }
func BenchmarkSlowStreamingDecode(b *testing.B)     { runAllSlowStreamingDecode(b) }

func runAllSlowStreamingDecode(tb testing.TB) {
	for _, td := range slowStreamingDecodeTestdata {
		for _, typeName := range []string{"Token", "Value"} {
			name := path.Join(td.name, typeName)
			runTestOrBench(tb, name, len64(td.data), func(tb testing.TB) {
				runSlowStreamingDecode(tb, typeName, td.data)
			})
		}
	}
}

// runSlowStreamingDecode tests a streaming Decoder operating on
// a slow io.Reader that only returns 1 byte at a time,
// which tends to exercise pathological behavior.
func runSlowStreamingDecode(t testing.TB, typeName string, data []byte) {
	if benchV1 {
		dec := jsonv1.NewDecoder(iotest.OneByteReader(bytes.NewReader(data)))
		switch typeName {
		case "Token":
			for dec.More() {
				if _, err := dec.Token(); err != nil {
					t.Fatalf("Decoder.Token error: %v", err)
				}
			}
		case "Value":
			var val jsonv1.RawMessage
			if err := dec.Decode(&val); err != nil {
				t.Fatalf("Decoder.Decode error: %v", err)
			}
		}
		return
	}

	dec := jsontext.NewDecoder(iotest.OneByteReader(bytes.NewReader(data)))
	switch typeName {
	case "Token":
		for dec.PeekKind() > 0 {
			if _, err := dec.ReadToken(); err != nil {
				t.Fatalf("Decoder.ReadToken error: %v", err)
			}
		}
	case "Value":
		if _, err := dec.ReadValue(); err != nil {
			t.Fatalf("Decoder.ReadValue error: %v", err)
		}
	}
}

func TestBenchmarkTextValue(t *testing.T) { runValue(t) }
func BenchmarkTextValue(b *testing.B)     { runValue(b) }

func runValue(tb testing.TB) {
	if testing.Short() {
		tb.Skip() // CitmCatalog is not loaded in short mode
	}
	var data []byte
	for _, ts := range jsontest.Data {
		if ts.Name == "CitmCatalog" {
			data = ts.Data()
		}
	}

	runTestOrBench(tb, "IsValid", len64(data), func(tb testing.TB) {
		jsontext.Value(data).IsValid()
	})

	methods := []struct {
		name   string
		format func(*jsontext.Value) error
	}{
		{"Compact", (*jsontext.Value).Compact},
		{"Indent", func(v *jsontext.Value) error { return v.Indent("\t", "    ") }},
		{"Canonicalize", (*jsontext.Value).Canonicalize},
	}

	var v jsontext.Value
	for _, method := range methods {
		runTestOrBench(tb, method.name, len64(data), func(tb testing.TB) {
			v = append(v[:0], data...) // reset with original input
			if err := method.format(&v); err != nil {
				tb.Errorf("jsontext.Value.%v error: %v", method.name, err)
			}
		})
		v = append(v[:0], data...)
		method.format(&v)
		runTestOrBench(tb, method.name+"/Noop", len64(data), func(tb testing.TB) {
			if err := method.format(&v); err != nil {
				tb.Errorf("jsontext.Value.%v error: %v", method.name, err)
			}
		})
	}
}

func runTestOrBench(tb testing.TB, name string, numBytes int64, run func(tb testing.TB)) {
	switch tb := tb.(type) {
	case *testing.T:
		tb.Run(name, func(t *testing.T) {
			run(t)
		})
	case *testing.B:
		tb.Run(name, func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()
			b.SetBytes(numBytes)
			for range b.N {
				run(b)
			}
		})
	}
}
