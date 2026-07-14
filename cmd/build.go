/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

type WasmDir struct {
	inputDir string
	outputDir string
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "It can build your solid project",
	RunE: buildServer,
}

func buildServer(cmd *cobra.Command, args []string) error {
	if err := npmBuild(); err != nil {
		return err
	}

	err := exportWasm()

	if err != nil {
		return err
	}

	appCmd := exec.Command("go", "build", ".")
	
	appCmd.Stdout = os.Stdout
	appCmd.Stderr = os.Stderr

	if err := appCmd.Run(); err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
