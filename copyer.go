package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
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
			_, _ = fmt.Fprintf(os.Stderr, "GOLINE parser failed:%s", err.Error())
			os.Exit(1)
		}
		fileLine = int(fl)
	}
	fmt.Println(fileLine)
	dir, err := os.Getwd()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "get working directory failed:%s", err.Error())
		os.Exit(1)
	}

	flag.Usage = Usage
	flag.Parse()

	srcPkg, srcName := parseSrcDstFlagName(*srcFlag)
	dstPkg, dstName := parseSrcDstFlagName(*dstFlag)
	gArg := &GeneratorArg{
		FileName: dir + "/" + fileName,
		Line:     fileLine,
		Src:      srcName,
		Dst:      dstName,
		SrcPkg:   srcPkg,
		DstPkg:   dstPkg,
	}
	err = generate(gArg)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "generate failed:%+v", err)
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

func parseSrcDstFlagName(s string) (string, string) {
	ss := strings.Split(s, ".")
	if len(ss) == 1 {
		return "", ss[0]
	} else {
		return ss[0], ss[1]
	}
}
