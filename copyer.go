package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

var (
	srcFlag = flag.String("src", "", "src type name")
	dstFlag = flag.String("dst", "", "dst type name")
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of copyer:\n")
	fmt.Fprintf(os.Stderr, "\tcopyer [flags] -src typename")
	fmt.Fprintf(os.Stderr, "\tcopyer [flags] -dst typename")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}
func main() {
	fileName := os.Getenv("GOFILE")
	var fileLine int
	if str := os.Getenv("GOLINE"); str != "" {
		fl, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "GOLINE parse failed:%s", err.Error())
			os.Exit(1)
		}
		fileLine = int(fl)
	}

	flag.Usage = Usage
	flag.Parse()
	err := generate(&GeneratorArg{
		FileName: fileName,
		Line:     fileLine,
		Src:      *srcFlag,
		Dst:      *dstFlag,
	})
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "generate failed:%s", err.Error())
		os.Exit(1)
	}

}

func generate(arg *GeneratorArg) error {
	g, err := NewGenerator(arg)
	if err != nil {
		return err
	}
	data, err := g.Generate()
	if err != nil {
		return err
	}
	err = g.OutPut(data)
	if err != nil {
		return err
	}
	return nil
}
