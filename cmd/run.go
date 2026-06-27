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
	Run: runServer,
}

func runServer(cmd *cobra.Command, args []string) {
	viteCmd := exec.Command("npm", "run", "build")

	viteCmd.Stdout = os.Stdout
	viteCmd.Stderr = os.Stderr

	err := viteCmd.Run()
	if err != nil {
		os.Exit(1)
	}

	appCmd := exec.Command("go", "run", ".", "--debug")
	
	appCmd.Stdout = os.Stdout
	appCmd.Stderr = os.Stderr

	err = appCmd.Run()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(runCmd)
}
