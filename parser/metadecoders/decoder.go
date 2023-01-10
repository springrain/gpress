// Copyright 2018 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metadecoders

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/niklasfasching/go-org/org"

	xml "github.com/clbanning/mxj/v2"
	toml "github.com/pelletier/go-toml/v2"
	"github.com/spf13/afero"
	"github.com/spf13/cast"
	jww "github.com/spf13/jwalterweatherman"
	yaml "gopkg.in/yaml.v2"
)

// Decoder provides some configuration options for the decoders.
type Decoder struct {
	// Delimiter is the field delimiter used in the CSV decoder. It defaults to ','.
	Delimiter rune

	// Comment, if not 0, is the comment character used in the CSV decoder. Lines beginning with the
	// Comment character without preceding whitespace are ignored.
	Comment rune
}

// OptionsKey is used in cache keys.
func (d Decoder) OptionsKey() string {
	var sb strings.Builder
	sb.WriteRune(d.Delimiter)
	sb.WriteRune(d.Comment)
	return sb.String()
}

// Default is a Decoder in its default configuration.
var Default = Decoder{
	Delimiter: ',',
}

// UnmarshalToMap will unmarshall data in format f into a new map. This is
// what's needed for Hugo's front matter decoding.
func (d Decoder) UnmarshalToMap(data []byte, f Format) (map[string]any, error) {
	m := make(map[string]any)
	if data == nil {
		return m, nil
	}

	err := d.UnmarshalTo(data, f, &m)

	return m, err
}

// UnmarshalFileToMap is the same as UnmarshalToMap, but reads the data from
// the given filename.
func (d Decoder) UnmarshalFileToMap(fs afero.Fs, filename string) (map[string]any, error) {
	format := FormatFromString(filename)
	if format == "" {
		return nil, fmt.Errorf("%q is not a valid configuration format", filename)
	}

	data, err := afero.ReadFile(fs, filename)
	if err != nil {
		return nil, err
	}
	return d.UnmarshalToMap(data, format)
}

// UnmarshalStringTo tries to unmarshal data to a new instance of type typ.
func (d Decoder) UnmarshalStringTo(data string, typ any) (any, error) {
	data = strings.TrimSpace(data)
	// We only check for the possible types in YAML, JSON and TOML.
	switch typ.(type) {
	case string:
		return data, nil
	case map[string]any:
		format := d.FormatFromContentString(data)
		return d.UnmarshalToMap([]byte(data), format)
	case []any:
		// A standalone slice. Let YAML handle it.
		return d.Unmarshal([]byte(data), YAML)
	case bool:
		return cast.ToBoolE(data)
	case int:
		return cast.ToIntE(data)
	case int64:
		return cast.ToInt64E(data)
	case float64:
		return cast.ToFloat64E(data)
	default:
		return nil, fmt.Errorf("unmarshal: %T not supported", typ)
	}
}

// Unmarshal will unmarshall data in format f into an interface{}.
// This is what's needed for Hugo's /data handling.
func (d Decoder) Unmarshal(data []byte, f Format) (any, error) {
	if len(data) == 0 {
		switch f {
		case CSV:
			return make([][]string, 0), nil
		default:
			return make(map[string]any), nil
		}
	}
	var v any
	err := d.UnmarshalTo(data, f, &v)

	return v, err
}

// UnmarshalTo unmarshals data in format f into v.
func (d Decoder) UnmarshalTo(data []byte, f Format, v any) error {
	var err error

	switch f {
	case ORG:
		err = d.unmarshalORG(data, v)
	case JSON:
		err = json.Unmarshal(data, v)
	case XML:
		var xmlRoot xml.Map
		xmlRoot, err = xml.NewMapXml(data)

		var xmlValue map[string]any
		if err == nil {
			xmlRootName, err := xmlRoot.Root()
			if err != nil {
				return err
			}
			xmlValue = xmlRoot[xmlRootName].(map[string]any)
		}

		switch v := v.(type) {
		case *map[string]any:
			*v = xmlValue
		case *any:
			*v = xmlValue
		}
	case TOML:
		err = toml.Unmarshal(data, v)
	case YAML:
		err = yaml.Unmarshal(data, v)
		if err != nil {
			return err
		}

		// To support boolean keys, the YAML package unmarshals maps to
		// map[interface{}]interface{}. Here we recurse through the result
		// and change all maps to map[string]interface{} like we would've
		// gotten from `json`.
		var ptr any
		switch v.(type) {
		case *map[string]any:
			ptr = *v.(*map[string]any)
		case *any:
			ptr = *v.(*any)
		default:
			// Not a map.
		}

		if ptr != nil {
			if mm, changed := stringifyMapKeys(ptr); changed {
				switch v.(type) {
				case *map[string]any:
					*v.(*map[string]any) = mm.(map[string]any)
				case *any:
					*v.(*any) = mm
				}
			}
		}
	case CSV:
		return d.unmarshalCSV(data, v)

	default:
		return fmt.Errorf("unmarshal of format %q is not supported", f)
	}

	if err == nil {
		return nil
	}

	return err
}

func (d Decoder) unmarshalCSV(data []byte, v any) error {
	r := csv.NewReader(bytes.NewReader(data))
	r.Comma = d.Delimiter
	r.Comment = d.Comment

	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	switch v.(type) {
	case *any:
		*v.(*any) = records
	default:
		return fmt.Errorf("CSV cannot be unmarshaled into %T", v)

	}

	return nil
}

func parseORGDate(s string) string {
	r := regexp.MustCompile(`[<\[](\d{4}-\d{2}-\d{2}) .*[>\]]`)
	if m := r.FindStringSubmatch(s); m != nil {
		return m[1]
	}
	return s
}

func (d Decoder) unmarshalORG(data []byte, v any) error {
	config := org.New()
	config.Log = jww.WARN
	document := config.Parse(bytes.NewReader(data), "")
	if document.Error != nil {
		return document.Error
	}
	frontMatter := make(map[string]any, len(document.BufferSettings))
	for k, v := range document.BufferSettings {
		k = strings.ToLower(k)
		if strings.HasSuffix(k, "[]") {
			frontMatter[k[:len(k)-2]] = strings.Fields(v)
		} else if k == "tags" || k == "categories" || k == "aliases" {
			jww.WARN.Printf("Please use '#+%s[]:' notation, automatic conversion is deprecated.", k)
			frontMatter[k] = strings.Fields(v)
		} else if k == "date" {
			frontMatter[k] = parseORGDate(v)
		} else {
			frontMatter[k] = v
		}
	}
	switch v.(type) {
	case *map[string]any:
		*v.(*map[string]any) = frontMatter
	default:
		*v.(*any) = frontMatter
	}
	return nil
}

// stringifyMapKeys recurses into in and changes all instances of
// map[interface{}]interface{} to map[string]interface{}. This is useful to
// work around the impedance mismatch between JSON and YAML unmarshaling that's
// described here: https://github.com/go-yaml/yaml/issues/139
//
// Inspired by https://github.com/stripe/stripe-mock, MIT licensed
func stringifyMapKeys(in any) (any, bool) {
	switch in := in.(type) {
	case []any:
		for i, v := range in {
			if vv, replaced := stringifyMapKeys(v); replaced {
				in[i] = vv
			}
		}
	case map[string]any:
		for k, v := range in {
			if vv, changed := stringifyMapKeys(v); changed {
				in[k] = vv
			}
		}
	case map[any]any:
		res := make(map[string]any)
		var (
			ok  bool
			err error
		)
		for k, v := range in {
			var ks string

			if ks, ok = k.(string); !ok {
				ks, err = cast.ToStringE(k)
				if err != nil {
					ks = fmt.Sprintf("%v", k)
				}
			}
			if vv, replaced := stringifyMapKeys(v); replaced {
				res[ks] = vv
			} else {
				res[ks] = v
			}
		}
		return res, true
	}

	return nil, false
}
