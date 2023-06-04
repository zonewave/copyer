package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "copyer",
	Short: "copyer is a tool to generate the copy code for golang",
	// TODO: add long description
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("copyer is starting...")
		cmdFlag, err := RootCmdFlagGet(cmd)
		if err != nil {
			logrus.Errorf("RootCmdFlagGet error:%v", err)
			return
		}
		env, err := NewEnv()
		if err != nil {
			logrus.Errorf("Env error:%v", err)
			return
		}
		err = LocalCopy(cmdFlag, env)
		if err != nil {
			logrus.Errorf("copyer error:%+v", err)
		}

		logrus.WithFields(logrus.Fields{
			"flags": cmdFlag,
		}).Info("copyer is end...")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	initLogrus()
	flags := rootCmd.PersistentFlags()
	flags.StringP("src", "s", "", "src type name")
	flags.StringP("dst", "d", "", "dst type name")
	flags.BoolP("print", "p", false, "print the copy code")

	rootCmd.AddCommand(outfileCmd)
	outfileCmd.Flags().StringP("out", "o", "", "out file name;default,it is copy_{{dst}}.go")
	outfileCmd.Flags().String("package", "", "package name;default, it is current packageName")
}

func initLogrus(opts ...func()) {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)
	for _, opt := range opts {
		opt()
	}
}
