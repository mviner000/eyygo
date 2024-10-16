package lib

import (
	"github.com/mviner000/eyygo/src/cmdlib"
	"github.com/spf13/cobra"
)

// GetRootCommand returns the root cobra command
func GetRootCommand() *cobra.Command {
	return rootCmd
}

// ExecuteRootCommand executes the root command
func ExecuteRootCommand() error {
	return rootCmd.Execute()
}

// RegisterCommand allows external users to register their own commands
func RegisterCommand(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}

// InitializeCommands sets up the initial command structure
func InitializeCommands() {
	cmdlib.RegisterInternalCommands(rootCmd)
}
