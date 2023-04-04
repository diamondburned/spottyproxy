package crudkv

// BasicStore is an untyped raw-bytes key-value store.
type BasicStore interface {
	Bucket(name string) BasicStore
	Tx(rw bool, f func(BasicTransaction) error) error
}

// BasicTransaction is an untyped raw-bytes transaction.
type BasicTransaction interface {
	Get(key string, v func([]byte) error) error
	Set(key string, value []byte) error
	Delete(key string) error
	ForEach(values bool, f func(string, []byte) error) error
}
