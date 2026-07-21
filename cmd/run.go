/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
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
	viteCmd := exec.Command("npm", "run", "build:watch")

	viteCmd.Stdout = os.Stdout
	viteCmd.Stderr = os.Stderr

	if err := viteCmd.Start(); err != nil {
		return err
	}

	defer func() {
        if viteCmd.Process != nil {
			if err := killProcessTree(viteCmd.Process.Pid); err != nil {
				fmt.Printf("warning: failed to kill process tree: %v\n", err)
			} else {
				fmt.Println("process tree killed")
			}
			_ = viteCmd.Process.Kill()
		}
    }()

	if err := exportWasm(); err != nil {
		return err
	}

	dirs, err := getSubDirNames("./cmd")

	if err != nil {
		return err
	}

	appCmd := exec.Command("go", "run", "./cmd/" + dirs[0], "--debug")
	
	appCmd.Stdout = os.Stdout
	appCmd.Stderr = os.Stderr

	if err := appCmd.Run(); err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(runCmd)
}
