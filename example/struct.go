package example

type Foo struct {
	Name   string
	Number int
	number int
	f      float32
}

type Boo struct {
	Name   string
	Number int
	number int
	f      float32
}

//go:generate ../bin/copyer -src=Foo -dst=Boo
