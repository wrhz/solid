/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

// modelCmd represents the model command
var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: createModel,
}

var orm string

func createModel(cmd *cobra.Command, args []string) error {
	if orm != "gorm" && orm != "xorm" && orm != "both" {
		return fmt.Errorf("%s", "The orm only can be gorm or xorm or both, are not " + orm)
	}

	modelName := strcase.ToCamel(args[0])
	
	if err := os.MkdirAll("./model", 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(filepath.Join("model", modelName + ".go"), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("The file is exist")
		}

		return err
	}
	defer f.Close()

	tag := "`"

	switch orm {
	case "gorm":
		tag += "gorm:\"primaryKey;autoIncrement\""
	case "xorm":
		tag += "xorm:\"'id' pk autoincr\""
	case "both":
		tag += "gorm:\"primaryKey;autoIncrement\" xorm:\"'id' pk autoincr\""
	}

	tag += "`"

	data := fmt.Sprintf(`package model

type %s struct {
	ID int64 %s

	// Write your struct members
}
`, modelName, tag)

	if _, err := f.Write([]byte(data)); err != nil {
		return err
	}

	return nil
}

func init() {
	modelCmd.Flags().StringVar(&orm, "orm", "both", "Choose model's orm(xorm or gorm or both)")

	createCmd.AddCommand(modelCmd)
}
