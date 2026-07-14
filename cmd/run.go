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

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "It can run your solid project and start the debug",
	RunE: runServer,
}

func runServer(cmd *cobra.Command, args []string) error {
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
		var buildCmd *exec.Cmd

		inputPath := "./" + filepath.Join("resource", "wasm", dir.inputDir)
		outputPath := filepath.Join(".", "dist", "resource", "wasm", dir.outputDir)

		buildCmd = exec.Command("go", "build", "-o", outputPath, inputPath)

		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stderr
		buildCmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")

		if err := buildCmd.Run(); err != nil {
			return err
		}
	}

	appCmd := exec.Command("go", "run", ".", "--debug")
	
	appCmd.Stdout = os.Stdout
	appCmd.Stderr = os.Stderr

	err = appCmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(runCmd)
}
