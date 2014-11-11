// Package stores provides a set of stores for streamtools
package stores

// a Store handles access to data that persists
type Store struct {
	Name      string
	BackingDB string
}

// create a new Store with a specified name
func NewStore(name string) *Store {
	return &Store{
		Name: name,
	}
}
