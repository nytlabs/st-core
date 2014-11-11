// Package stores provides a set of stores for streamtools
package stores

// a Store handles access to data that persists
type Store struct {
	Name      string
	BackingDB string
	QuitChan  chan bool
}

// create a new Store with a specified name
func NewStore(name string) *Store {
	return &Store{
		Name:     name,
		QuitChan: make(chan bool),
	}
}

func (s *Store) Stop() {
	s.QuitChan <- true
}
