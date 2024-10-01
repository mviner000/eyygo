package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var StartAppCmd = &cobra.Command{
	Use:   "startapp [appname]",
	Short: "Create a new application",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]
		fmt.Printf("Creating new app: %s\n", appName)
		// Implement app creation logic here
	},
}
