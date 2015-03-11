package core

import "errors"

func Copy(i interface{}) interface{} {
	switch t := i.(type) {
	case map[string]interface{}:
		n := make(map[string]interface{})
		for k, v := range t {
			n[k] = Copy(v)
		}
		return n
	case []interface{}:
		n := make([]interface{}, len(t), len(t))
		for i, v := range t {
			n[i] = Copy(v)
		}
		return n
	}
	return i
}

// this is the recursion called by mergeMap
func merge(A map[string]interface{}, i interface{}) interface{} {
	switch t := i.(type) {
	case map[string]interface{}:
		for k, v := range t {
			Ak, ok := A[k]
			if !ok {
				// if A doesn't contain this key, initialise
				switch v.(type) {
				case map[string]interface{}:
					// if the value is a map, initialise a new map to be merged
					Ak = make(map[string]interface{})
				default:
					// otherwise just merge into this A
					A[k] = v
					continue
				}
			}
			B, ok := Ak.(map[string]interface{})
			if !ok {
				// A[k] is not a map, so just overwrite
				A[k] = v
			}
			// if A contains the key, and the current value is also a map, recurse on that map
			A[k] = merge(B, v)
		}
		return A
	}
	// if a leaf is reached, just return it
	return i
}

// MergeMap recursively merges one map into another, overwriting base when necessary
func MergeMap(base, tomerge map[string]interface{}) (map[string]interface{}, error) {
	// call the recursion knowing the inputs are of the correct type
	result := merge(base, tomerge)
	// assert back to map
	out, ok := result.(map[string]interface{})
	if !ok {
		return nil, errors.New("failed merge")
	}
	return out, nil
}
