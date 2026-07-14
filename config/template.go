package config

import (
	"io"

	"github.com/wrhz/solid/template"
)

type TemplateConfigStruct struct {
	templateRender func(stdout io.Writer, name string, html string, args any) error
}

func (t *TemplateConfigStruct) SetTemplateRender(templateRender func(stdout io.Writer, name string, html string, args any) error) {
	t.templateRender = templateRender
}

func (t *TemplateConfigStruct) GetTemplateRender() func(stdout io.Writer, name string, html string, args any) error {
	return t.templateRender
}

func NewTemplateConfigStruct() *TemplateConfigStruct {
	return &TemplateConfigStruct{
		templateRender: template.TemplateRender,
	}
}