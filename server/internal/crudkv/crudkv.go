package crudkv

import (
	"errors"

	"github.com/fxamacker/cbor/v2"
)

// ErrNotFound is returned when a key is not found.
var ErrNotFound = errors.New("not found")

// Stop is a sentinel error that can be returned by a CRUD operation to
// indicate that the operation should stop. It is not an error in the
// traditional sense, but rather a signal to the caller that the operation
// should stop.
var Stop error = stop{}

type stop struct{}

func (stop) Error() string { return "stop" }

// Store is a typed key-value store.
type Store[V any] struct {
	store BasicStore
}

// Wrap wraps a BasicStore into a typed Store.
func Wrap[V any](store BasicStore, buckets ...string) Store[V] {
	if len(buckets) > 0 {
		for _, bucket := range buckets {
			store = store.Bucket(bucket)
		}
	}
	return Store[V]{store}
}

// Tx begins a transaction on the given store.
func (s Store[V]) Tx(rw bool, f func(Transaction[V]) error) error {
	return s.store.Tx(rw, func(txn BasicTransaction) error {
		return f(Transaction[V]{txn})
	})
}

// Get is a convenience method that gets a value from the store.
func (s Store[V]) Get(k string) (V, error) {
	var v V
	err := s.Tx(false, func(txn Transaction[V]) error {
		v2, err := txn.Get(k)
		v = v2
		return err
	})
	return v, err
}

// Transaction is a transaction on a key-value store.
type Transaction[V any] struct {
	txn BasicTransaction
}

func (tx Transaction[V]) Get(key string) (V, error) {
	var v V
	err := tx.txn.Get(key, func(b []byte) error {
		return cbor.Unmarshal(b, &v)
	})
	return v, err
}

func (tx Transaction[V]) Set(key string, value V) error {
	b, err := cbor.Marshal(value)
	if err != nil {
		return err
	}

	return tx.txn.Set(key, b)
}

func (tx Transaction[V]) Delete(key string) error {
	return tx.txn.Delete(key)
}

func (tx Transaction[V]) ForEach(f func(string, V) error) error {
	return tx.txn.ForEach(true, func(key string, b []byte) error {
		var v V
		if err := cbor.Unmarshal(b, &v); err != nil {
			return err
		}

		return f(key, v)
	})
}

func (tx Transaction[V]) ForEachLazy(f func(string, func() (V, error)) error) error {
	return tx.txn.ForEach(false, func(key string, _ []byte) error {
		return f(key, func() (V, error) {
			var v V
			err := tx.txn.Get(key, func(b []byte) error {
				return cbor.Unmarshal(b, &v)
			})
			return v, err
		})
	})
}
