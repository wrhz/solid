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

var routeType string

var createRouteCmd = &cobra.Command{
	Use:   "createRoute",
	Short: "You can use it to create routes",
	RunE: createRoute,
}

func createRoute(cmd *cobra.Command, args []string) error {
	if routeType != "main" && routeType != "common" {
		return fmt.Errorf("The route's type only can be main or common, are not %s", routeType)
	}

	routeName := strcase.ToCamel(args[0])
	varName := strcase.ToLowerCamel(args[0])

	if err := os.MkdirAll("./route", 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(filepath.Join("route", routeName + ".go"), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("The file is exist")
		}

		return err
	}
	defer f.Close()

	data := fmt.Sprintf(`package route

import "github.com/wrhz/Solid"

type %s struct {
	// Write your members
}

func New%s() *%s {
	// Write your code to new this struct

	return &%s{}
}

func (%s *%s) Init(r *solid.RouteStruct) {
	// Write your init code
}

func (%s *%s) RegisterRoute(r *solid.RouteStruct) {
	// Register your routes
}	

func (%s *%s) RegisterMiddleware(r *solid.RouteStruct) {
	// Register your moddlewares
}
`, routeName, routeName, routeName, routeName, varName, routeName, varName, routeName, varName, routeName)

	if routeType == "main" {
		data += fmt.Sprintf(`
func (%s *%s) ServerStart() {
	// Write your code, and it will run when the server starts.
}

func (%s *%s) ServerEnd() {
	// Write your code, and it will run when the server finishes.
}`, varName, routeName, varName, routeName)
	}

	if _, err := f.Write([]byte(data)); err != nil {
		return err
	}

	return nil
}

func init() {
	createRouteCmd.Flags().StringVar(&routeType, "type", "common", "Choose route's type(main or common)")

	rootCmd.AddCommand(createRouteCmd)
}
