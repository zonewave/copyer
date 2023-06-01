package cmd

import (
	"encoding/json"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func getString(set *pflag.FlagSet, name string) (string, error) {
	ret, err := set.GetString(name)
	if err != nil {
		return "", errors.Wrapf(err, "get %s string failed", name)
	}
	return ret, nil
}

type RootCmdFlag struct {
	Src   string
	Dst   string
	Print bool
}

func (r *RootCmdFlag) String() string {
	bs, _ := json.Marshal(r)
	return string(bs)
}

func RootCmdFlagGet(cmd *cobra.Command) (*RootCmdFlag, error) {
	flags := cmd.Flags()

	src, err := getString(flags, "src")
	if err != nil {
		return nil, err
	}
	dst, err := getString(flags, "dst")
	if err != nil {
		return nil, err
	}
	printOut, err := flags.GetBool("print")
	if err != nil {
		return nil, errors.Wrapf(err, "get print flag failed")
	}
	return &RootCmdFlag{
		Src:   src,
		Dst:   dst,
		Print: printOut,
	}, nil
}

type OutfileCmdFlag struct {
	Out string
	Pkg string
	*RootCmdFlag
}

func (r *OutfileCmdFlag) String() string {
	bs, _ := json.Marshal(r)
	return string(bs)
}
func OutfileCmdFlagGet(cmd *cobra.Command) (*OutfileCmdFlag, error) {
	flags := cmd.Flags()

	out, err := getString(flags, "out")
	if err != nil {
		return nil, err
	}
	pkg, err := getString(flags, "package")
	if err != nil {
		return nil, err
	}
	rootFlag, err := RootCmdFlagGet(cmd)
	if err != nil {
		return nil, err
	}

	return &OutfileCmdFlag{
		out,
		pkg,
		rootFlag,
	}, nil
}
