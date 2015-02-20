package core

import "errors"

func (j JSONType) MarshalJSON() ([]byte, error) {
	switch {
	case j == NUMBER:
		return []byte("Number"), nil
	case j == ARRAY:
		return []byte("Array"), nil
	case j == OBJECT:
		return []byte("Object"), nil
	case j == STRING:
		return []byte("String"), nil
	case j == BOOLEAN:
		return []byte("Boolean"), nil
	case j == ANY:
		return []byte("Any"), nil
	}
	return nil, errors.New("Unknown pin type")
}
