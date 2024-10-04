package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mviner000/eyymi/config"
	"github.com/spf13/cobra"
)

var DbmateCmd = &cobra.Command{
	Use:   "dbmate [command]",
	Short: "Run dbmate commands",
	Long:  `This command allows you to run dbmate operations like migrate, create, etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please provide a dbmate command (e.g., migrate, create, etc.)")
			return
		}

		// Get the absolute path of the project root
		projectRoot, err := filepath.Abs(config.ProjectRoot)
		if err != nil {
			fmt.Printf("Error getting absolute path: %v\n", err)
			os.Exit(1)
		}

		// Set the migrations directory
		migrationsDir := filepath.Join(projectRoot, "db", "migrations")

		dbmateCmd := exec.Command("dbmate", args...)
		dbmateCmd.Env = append(os.Environ(),
			fmt.Sprintf("DATABASE_URL=%s", config.GetDatabaseURLForDbmate()),
			fmt.Sprintf("MIGRATIONS_DIR=%s", migrationsDir),
		)
		dbmateCmd.Stdout = os.Stdout
		dbmateCmd.Stderr = os.Stderr

		err = dbmateCmd.Run()
		if err != nil {
			fmt.Printf("Error running dbmate: %v\n", err)
			os.Exit(1)
		}
	},
}
