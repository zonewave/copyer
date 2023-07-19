package xutil

import (
	"github.com/cockroachdb/errors"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/samber/mo"
)

func FlatMap[T, V any](m mo.Result[T], fn func(T) mo.Result[V]) mo.Result[V] {
	if m.IsError() {
		return mo.Err[V](m.Error())
	}
	return fn(m.MustGet())
}
func FlatMap2[T1, T2, V any](m1 mo.Result[T1], m2 mo.Result[T2], fn func(T1, T2) mo.Result[V]) mo.Result[V] {
	if m1.IsError() {
		return mo.Err[V](m1.Error())
	}
	if m2.IsError() {
		return mo.Err[V](m2.Error())
	}
	return fn(m1.MustGet(), m2.MustGet())
}

func Map[T, V any](m mo.Result[T], fn func(T) V) mo.Result[V] {
	return MapE(m, func(t T) (V, error) {
		return fn(t), nil
	})
}

func MapE[T, V any](m mo.Result[T], fn func(T) (V, error)) mo.Result[V] {
	if m.IsError() {
		return mo.Err[V](m.Error())
	}
	return mo.TupleToResult(fn(m.MustGet()))
}

func Map2[T1, T2, V any](m1 mo.Result[T1], m2 mo.Result[T2], fn func(T1, T2) V) mo.Result[V] {
	if m1.IsError() {
		return mo.Err[V](m1.Error())
	}
	if m2.IsError() {
		return mo.Err[V](m2.Error())
	}
	return mo.Ok(fn(m1.MustGet(), m2.MustGet()))
}
func Map2E[T1, T2, V any](m1 mo.Result[T1], m2 mo.Result[T2], fn func(T1, T2) (V, error)) mo.Result[V] {
	if m1.IsError() {
		return mo.Err[V](m1.Error())
	}
	if m2.IsError() {
		return mo.Err[V](m2.Error())
	}
	return mo.TupleToResult(fn(m1.MustGet(), m2.MustGet()))
}
func ChecksItems[K comparable, V any](items ...K) func(values map[K]V) (map[K]V, error) {
	return func(values map[K]V) (map[K]V, error) {
		if len(values) != len(items) {
			return nil, errors.Errorf("search:%v,but %v found", items, slice.Difference(items, maputil.Keys(values)))
		}
		return values, nil
	}
}
func MapWrap[T any](message string) func(err error) (T, error) {
	return func(err error) (v T, e error) {
		return v, errors.Wrap(err, message)
	}
}
func MapWrapf[T any](format string, args ...any) func(err error) (T, error) {

	return func(err error) (v T, e error) {
		return v, errors.Wrapf(err, format, args...)
	}
}
func OptionToResult[T any](m mo.Option[T], err error) mo.Result[T] {
	if m.IsPresent() {
		return mo.Ok[T](m.MustGet())
	} else {
		return mo.Err[T](err)
	}
}
func OptionResult[T any](fn func() (T, bool), err error) mo.Result[T] {
	v, ok := fn()
	if ok {
		return mo.Ok[T](v)
	} else {
		return mo.Err[T](err)
	}
}
