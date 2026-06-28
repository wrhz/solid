/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/iancoleman/strcase"
)

var orm string

var createModelCmd = &cobra.Command{
	Use:   "createModel",
	Short: "You can use it to create a database model for your solid project",
	Args:  cobra.MinimumNArgs(1),
	RunE: createModel,
}

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
	createModelCmd.Flags().StringVar(&orm, "orm", "both", "Choose model's orm(xorm or gorm or both)")

	rootCmd.AddCommand(createModelCmd)
}
