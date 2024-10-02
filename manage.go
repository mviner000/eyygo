package main

import (
	"fmt"
	"os"

	"github.com/mviner000/eyymi/cmd"
	"github.com/mviner000/eyymi/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "manage",
	Short: "Project management tool for your Go application",
}

func init() {
	rootCmd.AddCommand(cmd.ServerCmd)
	rootCmd.AddCommand(cmd.MigrateCmd)
	rootCmd.AddCommand(cmd.StartAppCmd)
	// Add more commands as needed
}

func main() {
	config.InitConfig()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}