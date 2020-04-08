package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	date    = "unknown"
	commit  = "local"
)

// execCmd represents the exec command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays version of cinit",
	Run: func(cmd *cobra.Command, args []string) {
		PrintVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func PrintVersion() {
	fmt.Printf("Version: %s\n", version)
	fmt.Printf("Commit: %s\n", commit)
	fmt.Printf("Date: %s\n", date)
}
