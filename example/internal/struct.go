package internal

import _ "strings"

type Foo struct {
	Name      string
	Number    int
	f32       float32
	Maps      map[int]int
	Slices    []int
	SlicesPtr []*int
}

type Bar struct {
	Name      string
	Number    int
	number    int
	Maps      map[int]int
	Slices    []int
	SlicesPtr []*int
}
