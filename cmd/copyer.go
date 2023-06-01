package cmd

import (
	"os"
	"strconv"
	"strings"

	"github.com/cockroachdb/errors"
	generate "github.com/zonewave/copyer/generate"
)

func LocalCopy(flag *RootCmdFlag) error {
	fileName := os.Getenv("GOFILE")
	var fileLine int
	if str := os.Getenv("GOLINE"); str != "" {
		fl, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return errors.Wrap(err, "go line parser failed:%s")
		}
		fileLine = int(fl)
	}
	dir, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "get working directory failed:%s")
	}
	srcPkg, srcName := parseSrcDstFlagName(flag.Src)
	dstPkg, dstName := parseSrcDstFlagName(flag.Dst)
	gArg := &generate.GeneratorArg{
		FileName: dir + "/" + fileName,
		Line:     fileLine,
		Src:      srcName,
		Dst:      dstName,
		SrcPkg:   srcPkg,
		DstPkg:   dstPkg,
	}
	err = generateCode(gArg)
	if err != nil {
		return errors.Wrap(err, "generate failed")
	}
	return nil
}

func generateCode(arg *generate.GeneratorArg) error {
	g, err := generate.NewGenerator(arg)
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
