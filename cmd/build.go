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
	viteCmd := exec.Command("npm", "run", "build")

	viteCmd.Stdout = os.Stdout
	viteCmd.Stderr = os.Stderr

	if err := viteCmd.Run(); err != nil {
		return err
	}

	if err := exportWasm(); err != nil {
		return err
	}

	dirs, err := getSubDirNames("./cmd")

	if err != nil {
		return err
	}

	appCmd := exec.Command("go", "build", "./cmd/" + dirs[0])
	
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
