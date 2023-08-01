package xmo

import (
	"github.com/cockroachdb/errors"
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

func Map3[T1, T2, T3, V any](m1 mo.Result[T1], m2 mo.Result[T2], m3 mo.Result[T3], fn func(T1, T2, T3) V) mo.Result[V] {
	if m1.IsError() {
		return mo.Err[V](m1.Error())
	}
	if m2.IsError() {
		return mo.Err[V](m2.Error())
	}
	if m3.IsError() {
		return mo.Err[V](m3.Error())
	}
	return mo.Ok(fn(m1.MustGet(), m2.MustGet(), m3.MustGet()))
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
