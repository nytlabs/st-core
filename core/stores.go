package core

type Store struct {
	Name      string
	BackingDB string
}

func NewStore(name string) *Store {
	return &Store{
		Name: name,
	}
}
