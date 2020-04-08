package cmd

import (
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"

	initApp "github.com/techniumlabs/cinit/pkg/app"
)

var cfgFile string
var app *initApp.App

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cinit",
	Short: "Init process for containers",
	Long: `Init process for containers which does the following
1. Proper Signal Forwarding
2. Orphaned Zombies Reaping
3. Fetch Secrets and expose it as Environment Variables`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		app, err = initApp.NewApp(cfgFile)
		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		app.RunInit(args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to execute: %s", err.Error())
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cinit.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
