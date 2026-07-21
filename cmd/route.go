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

// routeCmd represents the route command
var routeCmd = &cobra.Command{
	Use:   "route",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
	routeCmd.Flags().StringVar(&routeType, "type", "common", "Choose route's type(main or common)")

	createCmd.AddCommand(routeCmd)
}
