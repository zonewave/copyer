package xslice

import (
	"github.com/duke-git/lancet/v2/slice"
	"github.com/samber/mo"
	"golang.org/x/exp/constraints"
)

func FindBy[T any](s []T, predicate func(index int, item T) bool) mo.Option[T] {
	return mo.TupleToOption(slice.FindBy(s, predicate))
}

func FindByR[T any](s []T, predicate func(index int, item T) bool, err error) mo.Result[T] {
	ret, ok := slice.FindBy(s, predicate)
	if ok {
		return mo.Ok[T](ret)
	} else {
		return mo.Err[T](err)
	}
}

func SliceSetClassId[T any, K comparable, D constraints.Integer](s []T, id D, classify func(int, T) K, setId func(T, D)) ([]T, map[K]D) {
	idMap := make(map[K]D)
	for i, item := range s {
		key := classify(i, item)
		preId, ok := idMap[key]
		if !ok {
			idMap[key] = id
			setId(item, id)
			id++
		} else {
			setId(item, preId)
		}
	}
	return s, idMap
}
