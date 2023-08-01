package cmd

import (
	"os"
	"strconv"

	"github.com/cockroachdb/errors"
	"github.com/samber/mo"
	"github.com/zonewave/copyer/xutil/xmo"
)

type Env struct {
	GoLine    int
	GoFile    string
	GoPackage string
}

func newEnv(goLine int, goFile string, goPackage string) *Env {
	return &Env{GoLine: goLine, GoFile: goFile, GoPackage: goPackage}
}

func GoLineGet() mo.Result[int] {
	str := os.Getenv("GOLINE")
	if str != "" {
		val, err := strconv.ParseInt(str, 10, 64)
		return mo.TupleToResult(int(val), err).
			MapErr(xmo.MapWrap[int]("parse env GOLINE error"))
	} else {
		return mo.Ok(0)
	}
}

func EnvStringGet(name string) mo.Result[string] {
	val := os.Getenv("GOFILE")
	if val == "" {
		return mo.Err[string](errors.Newf("env %s is empty", name))
	}
	return mo.Ok(val)
}

func NewEnv() mo.Result[*Env] {
	return xmo.Map3(GoLineGet(), EnvStringGet("GOFILE"), EnvStringGet("GOPACKAGE"), newEnv)

}
