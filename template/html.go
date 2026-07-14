package template

import (
	"io"
	"text/template"
)

func TemplateRender(stdout io.Writer, name string, html string, args any) error {
	tmpl, err := template.New("demo").Parse(html)
    if err != nil {
        return err
    }

	err = tmpl.Execute(stdout, args)

	return err
}