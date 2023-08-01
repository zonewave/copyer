package xutil

import (
	"github.com/cockroachdb/errors"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/duke-git/lancet/v2/slice"
)

func ChecksItems[K comparable, V any](items ...K) func(values map[K]V) (map[K]V, error) {
	return func(values map[K]V) (map[K]V, error) {
		if len(values) != len(items) {
			return nil, errors.Errorf("search:%v,but %v found", items, slice.Difference(items, maputil.Keys(values)))
		}
		return values, nil
	}
}
