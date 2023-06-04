package generate

import (
	"github.com/zonewave/copyer/common"
	"golang.org/x/tools/go/packages"
)

type GeneratorArg struct {
	Action         common.ActionType
	GoFile         string
	GoLine         int
	OutFile        string
	OutLine        int
	SrcName        string
	SrcType        string
	SrcPkg         string
	DstName        string
	DstPkg         string
	DstType        string
	LoadConfigOpts []func(*packages.Config)
	Print          bool
}
