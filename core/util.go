package core

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
	default:
		return t
	}
	return nil
}
