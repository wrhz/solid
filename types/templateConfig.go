package types

import "io"

type ITemplateConfig interface {
	GetTemplateRender() func(stdout io.Writer, name string, html string, args any) error
}