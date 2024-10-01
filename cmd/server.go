package cmd

import (
	"github.com/mviner000/eyymi/core"
	"github.com/spf13/cobra"
)

var ServerCmd = &cobra.Command{
	Use:   "runserver",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		core.RunCommand()
	},
}
