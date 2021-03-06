package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays version of installed Rollout CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(GetCurrentVersion())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

// GetCurrentVersion of the CLI
func GetCurrentVersion() string {
	return "0.0.7" // ci-version-check
}
