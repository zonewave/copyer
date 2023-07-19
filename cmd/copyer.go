package cmd

import (
	"io"
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
	g := generate.NewGenerator(arg)
	bs := xutil.FlatMap(g, generate.ProduceCode)

	return xutil.FlatMap2(bs, mo.Ok(OutPutGet(arg)), generate.OutPut)

}
func OutPutGet(arg *generate.GeneratorArg) io.Writer {
	if arg.Print {
		return os.Stdout
	} else {
		return output.NewOutput(arg.OutFile, arg.OutLine)
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
