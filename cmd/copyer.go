package cmd

import (
	"os"
	"strings"

	"github.com/samber/mo"
	"github.com/zonewave/copyer/common"
	generate "github.com/zonewave/copyer/generate"
	"github.com/zonewave/copyer/output"
	"github.com/zonewave/copyer/xutil/xmo"
)

func LocalCopy(param *RootParam) error {
	env := param.Env
	flag := param.CmdFlag
	dir := mo.TupleToResult(os.Getwd()).
		MapErr(xmo.MapWrap[string]("get current dir error"))

	gArg := xmo.Map(dir, func(dir string) *generate.GeneratorArg {
		return &generate.GeneratorArg{
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
	})
	return xmo.FlatMap(
		gArg,
		generateCode,
	).Error()
}

func generateCode(arg *generate.GeneratorArg) mo.Result[bool] {

	outData := xmo.FlatMap(
		generate.NewGenerator(arg),
		generate.ProduceCode,
	)

	return xmo.FlatMap(
		outData,
		func(data []*output.LinesData) mo.Result[bool] {
			out := OutPutGet(arg)
			return generate.OutPut(data, out)
		},
	)

}
func OutPutGet(arg *generate.GeneratorArg) output.Writer {
	if arg.Print {
		return output.NewStdout()
	} else {
		return output.NewFile(arg.OutFile)
	}
}
func parseSrcDstFlagName(s string) (string, string) {
	ss := strings.Split(s, ".")
	if len(ss) == 1 {
		return "", ss[0]
	} else {
		return ss[0], ss[1]
	}
}
