package lib

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	conf "github.com/mviner000/eyygo/src"
	"github.com/mviner000/eyygo/src/cmdlib"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "manage",
	Short: "Project management tool for your Go application",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Load environment variables from .env file
		if err := godotenv.Load(); err != nil {
			fmt.Println("No .env file found")
		}
		conf.GetSettings()
	},
}

func init() {
	// Register internal commands
	cmdlib.RegisterInternalCommands(rootCmd)

	// Add custom commands here
	// Example:
	// rootCmd.AddCommand(customCmd)

	// Note for developers:
	// To add a new custom command, define your command using the cobra.Command struct.
	// Add your command here using rootCmd.AddCommand(yourCommand).
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
