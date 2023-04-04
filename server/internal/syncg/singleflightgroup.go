package syncg

import "golang.org/x/sync/singleflight"

type SingleflightGroup[V any] singleflight.Group

func (g *SingleflightGroup[V]) Do(key string, f func() (V, error)) (V, error) {
	iface, err, _ := (*singleflight.Group)(g).Do(key, func() (interface{}, error) {
		return f()
	})
	return iface.(V), err
}
