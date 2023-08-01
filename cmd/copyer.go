package cmd

import (
	"os"
	"strings"

	"github.com/samber/mo"
	"github.com/zonewave/copyer/common"
	generate "github.com/zonewave/copyer/generate"
	"github.com/zonewave/copyer/output"
	"github.com/zonewave/copyer/xutil"
)

func LocalCopy(flag *RootCmdFlag, env *Env) mo.Result[bool] {
	dir := mo.TupleToResult(os.Getwd()).
		MapErr(xutil.MapWrap[string]("get current dir error"))

	gArg := xutil.Map(dir, func(dir string) *generate.GeneratorArg {
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
	return xutil.FlatMap(gArg, generateCode)
}

func generateCode(arg *generate.GeneratorArg) mo.Result[bool] {

	bs := xutil.FlatMap(generate.NewGenerator(arg), generate.ProduceCode)

	return xutil.FlatMap2(bs, mo.Ok(OutPutGet(arg)), func(bs []byte, out output.Writer) mo.Result[bool] {
		return generate.OutPut(arg.GoLine, bs, out)
	})

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
