/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// protoCmd represents the proto command
var protoCmd = &cobra.Command{
	Use:   "proto",
	Short: "To build proto files",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: buildProto,
}

func getProtoFiles() ([]string, error) {
    var protoFiles []string

	root := "./proto"

    err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
        if err != nil {
            return err
        }

        if !d.IsDir() && strings.HasSuffix(d.Name(), ".proto") {
            protoFiles = append(protoFiles, path)
        }
        return nil
    })
    if err != nil {
        return nil, fmt.Errorf("Walk dir error: %w", err)
    }
    return protoFiles, nil
}

func buildProto(cmd *cobra.Command, args []string) error {
	protoFiles, err := getProtoFiles()

	if err != nil {
		return err
	}

	for _, filePath := range protoFiles {
		protoCmd := exec.Command("protoc",
			"--go_out=.",
			"--go_opt=paths=source_relative",
			"--go-grpc_out=.",
			"--go-grpc_opt=paths=source_relative",
			filePath,
		)

		protoCmd.Stdout = os.Stdout
		protoCmd.Stderr = os.Stderr

		if err := protoCmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(protoCmd)
}
