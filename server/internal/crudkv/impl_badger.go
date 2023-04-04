package crudkv

import (
	"errors"

	"github.com/dgraph-io/badger/v4"
)

type badgerStore struct {
	db     *badger.DB
	prefix string
}

// NewBadgerStore creates a new BadgerDB-backed store.
func NewBadgerStore(db *badger.DB, prefix string) BasicStore {
	return &badgerStore{db, prefix}
}

func (s *badgerStore) Bucket(name string) BasicStore {
	return &badgerStore{s.db, s.prefix + name + "\x00"}
}

func (s *badgerStore) Tx(rw bool, f func(BasicTransaction) error) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return f(&badgerTxn{txn, s.prefix})
	})
}

type badgerTxn struct {
	txn    *badger.Txn
	prefix string
}

var _ BasicTransaction = (*badgerTxn)(nil)

func (t *badgerTxn) Get(key string, v func([]byte) error) error {
	k := []byte(t.prefix + key)

	item, err := t.txn.Get(k)
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return ErrNotFound
		}
		return err
	}

	return item.Value(v)
}

func (t badgerTxn) Set(key string, value []byte) error {
	k := []byte(t.prefix + key)
	return t.txn.Set(k, value)
}

func (t badgerTxn) Delete(key string) error {
	k := []byte(t.prefix + key)
	return t.txn.Delete(k)
}

func (t badgerTxn) ForEach(values bool, f func(string, []byte) error) error {
	iter := t.txn.NewIterator(badger.IteratorOptions{
		PrefetchValues: values,
		Prefix:         []byte(t.prefix),
	})
	defer iter.Close()

	for iter.Rewind(); iter.Valid(); iter.Next() {
		item := iter.Item()
		key := string(item.Key())[len(t.prefix):]

		var value []byte
		if values {
			if err := item.Value(func(bvalue []byte) error {
				value = bvalue
				return nil
			}); err != nil {
				return err
			}
		}

		if err := f(key, value); err != nil {
			if err == Stop {
				return nil
			}
			return err
		}
	}

	return nil
}
