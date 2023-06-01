package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "copyer",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("copyer is starting...")
		src, _ := cmd.Flags().GetString("src")
		dst, _ := cmd.Flags().GetString("dst")

		err := Copyer(src, dst)
		if err != nil {
			logrus.Errorf("copyer error:%+v", err)
		}

		logrus.WithFields(logrus.Fields{
			"src": src,
			"dst": dst,
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
	flags := rootCmd.Flags()
	flags.StringP("src", "s", "", "src type name")
	flags.StringP("dst", "d", "", "dst type name")
	flags.StringP("output", "o", "", `
		output file name, default is copy_[src]_[dst].go;
		if output=="local" ，output to on the line of the file where the gender command resides`)

}
func initLogrus(opts ...func()) {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)
	for _, opt := range opts {
		opt()
	}
}
