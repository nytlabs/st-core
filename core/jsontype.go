package core

import "errors"

func (j JSONType) MarshalJSON() ([]byte, error) {
	switch j {
	case NUMBER:
		return []byte(`"number"`), nil
	case ARRAY:
		return []byte(`"array"`), nil
	case OBJECT:
		return []byte(`"object"`), nil
	case STRING:
		return []byte(`"string"`), nil
	case BOOLEAN:
		return []byte(`"boolean"`), nil
	case NULL:
		return []byte(`"null"`), nil
	case ANY:
		return []byte(`"any"`), nil
	}
	return nil, errors.New("Unknown pin type")
}
