package core

func GetSources() map[string]SourceSpec {
	sources := []SourceSpec{
		NSQInterface(),
		KeyValueStore(),
		ValueStore(),
		PriorityQueueStore(),
		ListStore(),
		WebsocketClient(),
	}

	library := make(map[string]SourceSpec)

	for _, s := range sources {
		library[s.Name] = s
	}

	return library
}

func isStore(sourceType SourceType) bool {
	_, ok := map[SourceType]bool{
		KEY_VALUE:       true,
		LIST:            true,
		VALUE_PRIMITIVE: true,
	}[sourceType]
	return ok
}

func isInterface(sourceType SourceType) bool {
	_, ok := map[SourceType]bool{
		NSQCLIENT: true,
		WSCLIENT:  true,
	}[sourceType]
	return ok
}
