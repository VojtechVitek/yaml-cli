package cli

import (
	"fmt"
)

// gopkg.in/yaml.v3 tries to unmarshal into map[string]interface{} whenever it can,
// but if finds YAML anchor, it falls back to generic map[interface{}]interface{},
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
