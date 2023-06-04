package cmd

import (
	"os"
	"strconv"

	"github.com/cockroachdb/errors"
)

type Env struct {
	GoLine    int
	GoFile    string
	GoPackage string
}

func NewEnv() (*Env, error) {
	var fileLine int
	if str := os.Getenv("GOLINE"); str != "" {
		fl, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, "go line parser failed:%s")
		}
		fileLine = int(fl)
	}
	return &Env{
		GoLine:    fileLine,
		GoFile:    os.Getenv("GOFILE"),
		GoPackage: os.Getenv("GOPACKAGE"),
	}, nil
}
