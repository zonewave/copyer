package cmd

import (
	"github.com/samber/mo"
	"github.com/spf13/cobra"
	"github.com/zonewave/copyer/xutil/xmo"
)

type RootParam struct {
	Env     *Env
	CmdFlag *RootCmdFlag
}

func NewRootParamCtr(env *Env, cmdFlag *RootCmdFlag) *RootParam {
	return &RootParam{Env: env, CmdFlag: cmdFlag}
}

func NewRootParam(cmd *cobra.Command) mo.Result[*RootParam] {
	return xmo.Map2(NewEnv(), RootCmdFlagGet(cmd), NewRootParamCtr)
}
