package cmd

import (
	"os"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/zonewave/copyer/common"
	generate "github.com/zonewave/copyer/generate"
)

func LocalCopy(flag *RootCmdFlag, env *Env) error {

	dir, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "get working directory failed:%s")
	}

	gArg := &generate.GeneratorArg{
		Action:         common.Local,
		GoFile:         dir + "/" + env.GoFile,
		GoLine:         env.GoLine,
		GoPkg:          env.GoPackage,
		OutFile:        dir + "/" + env.GoFile,
		OutLine:        env.GoLine,
		SrcName:        "src",
		SrcType:        flag.SrcType,
		SrcPkg:         flag.SrcPkg,
		DstName:        "dst",
		DstPkg:         flag.DstPkg,
		DstType:        flag.DstType,
		LoadConfigOpts: nil,
		Print:          false,
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
