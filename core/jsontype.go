package core

import "errors"

func (j *JSONType) UnmarshalJSON(data []byte) error {
	jsonType := string(data)
	switch jsonType {
	case `"number"`:
		*j = JSONType(NUMBER)
	case `"array"`:
		*j = JSONType(ARRAY)
	case `"object"`:
		*j = JSONType(OBJECT)
	case `"string"`:
		*j = JSONType(STRING)
	case `"boolean"`:
		*j = JSONType(BOOLEAN)
	case `"writer"`:
		*j = JSONType(WRITER)
	case `"null"`:
		*j = JSONType(NULL)
	case `"any"`:
		*j = JSONType(ANY)
	default:
		return errors.New("Error unmarshalling JSONType")
	}
	return nil
}

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
	case WRITER:
		return []byte(`"writer"`), nil
	case NULL:
		return []byte(`"null"`), nil
	case ANY:
		return []byte(`"any"`), nil
	}
	return nil, errors.New("Unknown pin type")
}
