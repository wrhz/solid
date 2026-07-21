package config

import "io"

type ItemplateRender interface {
	New(filename string) error
	Define(key string, value any)
	Render(name string, stdout io.Writer, args any) error
}

type ITemplateConfig interface {
	GetTemplateRender() ItemplateRender
}