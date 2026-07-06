/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "It can run your solid project and start the debug",
	RunE: runServer,
}

func runServer(cmd *cobra.Command, args []string) error {
	if err := npmBuild(); err != nil {
		return err
	}

	appCmd := exec.Command("go", "run", ".", "--debug")
	
	appCmd.Stdout = os.Stdout
	appCmd.Stderr = os.Stderr

	err := appCmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(runCmd)
}
