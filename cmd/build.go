/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "It can build your solid project",
	Run: buildServer,
}

func buildServer(cmd *cobra.Command, args []string) {
	viteCmd := exec.Command("npm", "run", "build")

	viteCmd.Stdout = os.Stdout
	viteCmd.Stderr = os.Stderr

	err := viteCmd.Run()
	if err != nil {
		os.Exit(1)
	}

	appCmd := exec.Command("go", "build", ".")
	
	appCmd.Stdout = os.Stdout
	appCmd.Stderr = os.Stderr

	err = appCmd.Run()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
