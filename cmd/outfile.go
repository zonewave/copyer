package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var outfileCmd = &cobra.Command{
	Use:   "outfile",
	Short: "copyer is a tool to generate the copy code for golang",
	Long:  `copyer is a tool to generate the copy code for golang.`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("copyer outfile is starting...")
		cmdFlag, err := OutfileCmdFlagGet(cmd).Get()
		if err != nil {
			logrus.Errorf("RootCmdFlagGet error:%v", err)
			return
		}

		logrus.WithFields(logrus.Fields{
			"cmdFlag": cmdFlag,
		}).Info("copyer is end...")
	},
}
