package store

type Store interface {
	List() ([][]byte, error)
	Get(key string) ([]byte, error)
	Apply(key string, value []byte) error
	Delete(key string) error
}
