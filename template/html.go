package template

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"
)

type HTMLTemplateRender struct {
	tmpls map[string]*template.Template
	defines template.FuncMap
}

func (h *HTMLTemplateRender) New(filename string) error {
	var html, err = os.ReadFile(filepath.Join(".", "resource", "views", filename))

	if err != nil {
		return err
	}

	tmpl, err := template.New(filename).Funcs(h.defines).Parse(string(html))

	if err != nil {
		return err
	}

	if h.tmpls == nil {
		h.tmpls = map[string]*template.Template{}
	}

	h.tmpls[filename] = tmpl

	return nil
}

func (h *HTMLTemplateRender) Render(name string, stdout io.Writer, args any) error {
	tmpl, ok := h.tmpls[name]

	if !ok {
		return fmt.Errorf("The %s is not in views", name)
	}

	return tmpl.Execute(stdout, args)
}

func (h *HTMLTemplateRender) Define(key string, value any) {
	if h.defines == nil {
		h.defines = map[string]any{}
	}

	h.defines[key] = func() any {
		return value
	}
}