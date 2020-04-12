package test

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// CompareJSON returns whether the given strings are equivalent as JSON.
// If needed, optional parameters can be excluded from the comparison.
func CompareJSON(got, want string, t *testing.T, ignoreFields ...string) string {
	var gotJSON, wantJSON interface{}

	if err := json.Unmarshal([]byte(want), &wantJSON); err != nil {
		return cmp.Diff(got, want)
	}
	if err := json.Unmarshal([]byte(got), &gotJSON); err != nil {
		t.Fatalf("failed to unmarshall got JSON: %s", got)
	}

	switchNRemove(gotJSON, ignoreFields)
	switchNRemove(wantJSON, ignoreFields)

	return cmp.Diff(gotJSON, wantJSON)
}

func switchNRemove(JSON interface{}, ignoreFields []string) {

	switch t := JSON.(type) {
	case map[string]interface{}:
		removeKeys(t, ignoreFields...)
	case []interface{}:
		for _, i := range t {
			if m, ok := i.(map[string]interface{}); ok {
				removeKeys(m, ignoreFields...)
			}
		}
	}
}

func removeKeys(m map[string]interface{}, keys ...string) {
	for _, k := range keys {
		delete(m, k)
	}

	for _, v := range m {
		switchNRemove(v, keys)
	}
}
