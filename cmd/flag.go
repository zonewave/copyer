package cmd

import (
	"encoding/json"

	"github.com/cockroachdb/errors"
	"github.com/samber/mo"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/zonewave/copyer/xutil/xmo"
)

func FlagsStringGet(set *pflag.FlagSet) func(name string) mo.Result[string] {
	return func(name string) mo.Result[string] {
		ret, err := set.GetString(name)
		return mo.TupleToResult(ret, errors.Wrapf(err, "get %s string failed", name))
	}
}
func FlagsBoolGet(set *pflag.FlagSet) func(name string) mo.Result[bool] {
	return func(name string) mo.Result[bool] {
		ret, err := set.GetBool(name)
		return mo.TupleToResult(ret, errors.Wrapf(err, "get %s bool failed", name))
	}
}

type RootCmdFlag struct {
	Src             string
	Dst             string
	SrcPkg, SrcType string
	DstPkg, DstType string
	Print           bool
}

func (r *RootCmdFlag) String() string {
	bs, _ := json.Marshal(r)
	return string(bs)
}

func RootCmdFlagGet(cmd *cobra.Command) mo.Result[*RootCmdFlag] {
	flags := cmd.Flags()
	flagStrings := FlagsStringGet(flags)
	return xmo.Map3(flagStrings("src"), flagStrings("dst"), FlagsBoolGet(flags)("print"), func(src, dst string, print bool) *RootCmdFlag {
		srcPkg, srcType := parseSrcDstFlagName(src)
		dstPkg, dstType := parseSrcDstFlagName(dst)
		return &RootCmdFlag{
			Src:     src,
			Dst:     dst,
			SrcPkg:  srcPkg,
			SrcType: srcType,
			DstPkg:  dstPkg,
			DstType: dstType,
			Print:   print,
		}
	})

}

type OutfileCmdFlag struct {
	Out string
	Pkg string
	*RootCmdFlag
}

func NewOutfileCmdFlag(out string, pkg string, rootCmdFlag *RootCmdFlag) *OutfileCmdFlag {
	return &OutfileCmdFlag{Out: out, Pkg: pkg, RootCmdFlag: rootCmdFlag}
}

func (r *OutfileCmdFlag) String() string {
	bs, _ := json.Marshal(r)
	return string(bs)
}
func OutfileCmdFlagGet(cmd *cobra.Command) mo.Result[*OutfileCmdFlag] {
	flags := cmd.Flags()
	flagStrings := FlagsStringGet(flags)
	return xmo.Map3(flagStrings("out"), flagStrings("package"), RootCmdFlagGet(cmd), NewOutfileCmdFlag)
}
