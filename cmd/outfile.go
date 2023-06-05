package cmd

import (
	"os"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/zonewave/copyer/common"
	"github.com/zonewave/copyer/generate"
)

var outfileCmd = &cobra.Command{
	Use:   "outfile",
	Short: "copyer is a tool to generate the copy code for golang",
	Long:  `copyer is a tool to generate the copy code for golang.`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("copyer outfile is starting...")
		cmdFlag, err := OutfileCmdFlagGet(cmd)
		if err != nil {
			logrus.Errorf("RootCmdFlagGet error:%v", err)
			return
		}
		env, err := NewEnv()
		if err != nil {
			logrus.Errorf("Env error:%v", err)
			return
		}
		err = OutFileCopy(cmdFlag, env)
		if err != nil {
			logrus.Errorf("copyer error:%+v", err)
		}
		logrus.WithFields(logrus.Fields{
			"cmdFlag": cmdFlag,
			"env":     env,
		}).Info("copyer is end...")
	},
}

func OutFileCopy(flag *OutfileCmdFlag, env *Env) error {

	dir, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "get working directory failed:%s")
	}

	gArg := &generate.GeneratorArg{
		Action:         common.Outfile,
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
