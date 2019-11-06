package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
	yamlv3 "gopkg.in/yaml.v3"
)

func yamlToJSON(w io.Writer, r io.Reader) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return errors.Wrap(err, "failed to read YAML data")
	}
	var data map[interface{}]interface{}
	if err := yamlv3.Unmarshal(b, &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal YAML data")
	}

	jsonCompatibleData := convertMap(data)

	// TODO: Marshal directly into w.
	b, err = json.MarshalIndent(jsonCompatibleData, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal to JSON")
	}
	_, err = w.Write(b)
	return err
}

func jsonToYAML(w io.Writer, r io.Reader) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return errors.Wrap(err, "failed to read JSON data")
	}
	var data map[string]interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal JSON data")
	}

	// TODO: Marshal directly into w.
	b, err = yamlv3.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal to YAML")
	}
	_, err = w.Write(b)
	return err
}

// gopkg.in/yaml.v3 tries to unmarshal into map[string]interface{} whenever it can,
// but if it finds YAML anchor, it falls back to generic map[interface{}]interface{},
// which in turn can't be handled by encoding/json. Thus we need to convert by hand.
func convertMap(genericMap map[interface{}]interface{}) map[string]interface{} {
	stringMap := make(map[string]interface{}, len(genericMap))
	for k, v := range genericMap {
		switch v := v.(type) {
		case map[interface{}]interface{}:
			stringMap[fmt.Sprint(k)] = convertMap(v)
		case []interface{}:
			for i := 0; i < len(v); i++ {
				switch v2 := v[i].(type) {
				case map[interface{}]interface{}:
					v[i] = convertMap(v2)
				}
			}
			stringMap[fmt.Sprint(k)] = v
		default:
			stringMap[fmt.Sprint(k)] = v
		}
	}
	return stringMap
}
