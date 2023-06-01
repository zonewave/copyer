package example

import (
	"strings"

	"github.com/zonewave/copyer/example/internal"
	in2 "github.com/zonewave/copyer/example/internal2"
)

type Foo struct {
	Name      string
	Number    int
	f32       float32
	Maps      map[int]int
	Slices    []int
	SlicesPtr []*int
	Foo       *Foo
}
type Bar struct {
	Name      string
	Number    int
	number    int
	Maps      map[int]int
	Slices    []int
	SlicesPtr []*int
	Foo       *Foo
}

func CopyFooToBar(src *Foo, dst *Bar) {
	dst.Foo = src.Foo
	dst.Maps = src.Maps
	dst.Name = src.Name
	dst.Number = src.Number
	dst.Slices = src.Slices
	dst.SlicesPtr = src.SlicesPtr
}
func CopyFooToBar(src *Foo, dst *Bar) {
	dst.Foo = src.Foo
	dst.Maps = src.Maps
	dst.Name = src.Name
	dst.Number = src.Number
	dst.Slices = src.Slices
	dst.SlicesPtr = src.SlicesPtr
}
//go:generate copyer local --src=Foo --dst=Bar

//go:generate copyer outfile -s internal.Foo -d Bar
var _ in2.Foo
var _ internal.Foo
var _ strings.Builder