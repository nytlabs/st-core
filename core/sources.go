package core

/*var SourceNames = map[sourceType]string{
	NONE:      "None",
	KEY_VALUE: "Key-Value",
	STREAM:    "Stream",
	PRIORITY:  "Priority Queue",
	ARRAY:     "Array",
	VALUE:     "Value",
}*/

type NewSource func() Source

func GetSources() map[string]NewSource {
	return map[string]NewSource{
		"keyvalue": NewKeyValue,
		"stream":   NewStream,
	}
}
