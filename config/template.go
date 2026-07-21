package config

import (
	"io"
	"os"
	"path/filepath"

	"github.com/wrhz/solid/template"
	"github.com/wrhz/solid/types/config"
)

type EmptyTemplateRender struct {
	html []byte
}

func (e *EmptyTemplateRender) New(filename string) error {
	var html, err = os.ReadFile(filepath.Join(".", "resource", "views", filename))

	if err != nil {
		return err
	}

	e.html = html

	return nil
}

func (e *EmptyTemplateRender) Render(_ string, stdout io.Writer, _ any) error {
	_, err := stdout.Write(e.html)

	return err
}

func (*EmptyTemplateRender) Define(_ string, _ any) {}

type TemplateConfigStruct struct {
	templateRender config.ItemplateRender
}

func (t *TemplateConfigStruct) SetTemplateRender(templateRender config.ItemplateRender) {
	t.templateRender = templateRender
}

func (t *TemplateConfigStruct) GetTemplateRender() config.ItemplateRender {
	if t.templateRender == nil {
		t.templateRender = &EmptyTemplateRender{}
	}

	return t.templateRender
}

func (t *TemplateConfigStruct) Define(key string, value any) {
	if t.templateRender != nil {
		t.templateRender.Define(key, value)
	}
}

func NewTemplateConfigStruct() *TemplateConfigStruct {
	return &TemplateConfigStruct{
		templateRender: &template.HTMLTemplateRender{},
	}
}