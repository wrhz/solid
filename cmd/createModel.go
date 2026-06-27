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
	Run: createModel,
}

func createModel(cmd *cobra.Command, args []string) {
	if orm != "gorm" && orm != "xorm" && orm != "both" {
		fmt.Println("The orm only can be gorm or xorm or both, are not " + orm)

		return
	}

	modelName := strcase.ToCamel(args[0])
	
	if err := os.MkdirAll("./model", 0755); err != nil {
		fmt.Println("Create directory error:", err)

		return
	}

	f, err := os.OpenFile(filepath.Join("model", modelName + ".go"), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			fmt.Println("The file is exist")

			return
		}
		fmt.Println("Create file error: %w", err)

		return
	}
	defer f.Close()

	tag := "`"

	if orm == "gorm" || orm == "both" {
		tag += "gorm:\"primaryKey;autoIncrement\" "
	}

	if orm == "xorm" || orm == "both" {
		tag += "xorm:\"'id' pk autoincr\""
	}

	tag += "`"

	data := fmt.Sprintf(`package model

type %s struct {
	ID int64 %s

	// Write your struct members
}`, modelName, tag)

	if _, err := f.Write([]byte(data)); err != nil {
		fmt.Println("Write file error: %w", err)
	}
}

func init() {
	createModelCmd.Flags().StringVar(&orm, "orm", "both", "Choose model's orm(xorm or gorm or both)")

	rootCmd.AddCommand(createModelCmd)
}
