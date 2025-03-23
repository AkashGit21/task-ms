package main

import (
	"os"

	"github.com/AkashGit21/task-ms/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "root",
	Short: "Root command of the project",
}

func init() {
	runTaskCmd := &cobra.Command{
		Use:   "task",
		Short: "Starts running the application server for task-service",
		Run: func(cmd *cobra.Command, args []string) {
			taskSrv, err := NewTaskV1Server()
			if err != nil {
				utils.ErrorLog("Error getting new server:", err)
				return
			}

			StartServer(taskSrv)
		},
	}

	runAuthnCmd := &cobra.Command{
		Use:   "authn",
		Short: "Starts running the application server for task-service",
		Run: func(cmd *cobra.Command, args []string) {
			authnSrv, err := NewAuthnServer()
			if err != nil {
				utils.ErrorLog("Error getting new server:", err)
				return
			}

			StartServer(authnSrv)
		},
	}

	rootCmd.AddCommand(runTaskCmd)
	rootCmd.AddCommand(runAuthnCmd)
}

func main() {
	Execute()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		utils.ErrorLog("Could not execute the application", err)
		os.Exit(1)
	}
}
