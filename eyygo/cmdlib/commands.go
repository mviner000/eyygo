package cmdlib

import (
	"github.com/mviner000/eyymi/eyygo/cmd"
	"github.com/spf13/cobra"
)

// RegisterInternalCommands adds internal commands to the provided root command.
func RegisterInternalCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(cmd.ServerCmd)
	rootCmd.AddCommand(cmd.StartAppCmd)
	rootCmd.AddCommand(cmd.MakeMigrationCmd)
	rootCmd.AddCommand(cmd.MigrateCmd)
	rootCmd.AddCommand(cmd.ShowRoutesCmd)
	rootCmd.AddCommand(cmd.CreateSuperuserCmd)
	rootCmd.AddCommand(cmd.MigratorCmd)
}
