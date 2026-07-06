/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"os/exec"
	"path/filepath"

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
	data, err := readWasmExecJS()

	if err != nil {
		return err
	}

	data = append(data, "\nexport default Go;"...)

	if err = os.WriteFile(filepath.Join("resource", "lib", "wasm_exec.js"), data, 0644); err != nil {
		return err
	}
	
	if err := npmBuild(); err != nil {
		return err
	}

	dirs, err := getDirs()

	if err != nil {
		return err
	}

	for _, dir := range dirs {
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(".", "dist", "resource", "wasm", dir.outputDir), "./" + filepath.Join("resource", "wasm", dir.inputDir))
	
		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stderr
		buildCmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")

		if err := buildCmd.Run(); err != nil {
			return err
		}
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
