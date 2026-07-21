/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "To test your Solid Project",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: TestProject,
}

func TestProject(cmd *cobra.Command, args []string) error {
	path, err := os.Getwd()

	if err != nil {
		return err
	}

	testCmd := exec.Command("go", "test", "-v", "./...")

	testCmd.Stdout = os.Stdout
	testCmd.Stderr = os.Stderr
	testCmd.Env = append(os.Environ(), "TEST_WORKDIR=" + path)

	return testCmd.Run()
}

func init() {
	rootCmd.AddCommand(testCmd)
}
