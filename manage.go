package main

import (
	"fmt"
	"os"

	"github.com/mviner000/eyymi/config"
	"github.com/mviner000/eyymi/eyygo/cmd"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "manage",
	Short: "Project management tool for your Go application",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.InitConfig()
	},
}

func init() {
	rootCmd.AddCommand(cmd.ServerCmd)
	rootCmd.AddCommand(cmd.StartAppCmd)
	rootCmd.AddCommand(cmd.DbmateCmd)
	rootCmd.AddCommand(cmd.MakeMigrationsCmd)
	rootCmd.AddCommand(cmd.MigrateCmd)
	// Add more commands as needed
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
