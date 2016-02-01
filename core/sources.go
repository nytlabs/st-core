package core

func GetSources() map[string]SourceSpec {
	sources := []SourceSpec{
		NSQConsumerInterface(),
		KeyValueStore(),
		ValueStore(),
		PriorityQueueStore(),
		ListStore(),
		WebsocketClient(),
		StdinInterface(),
	}

	library := make(map[string]SourceSpec)

	for _, s := range sources {
		library[s.Name] = s
	}

	return library
}
