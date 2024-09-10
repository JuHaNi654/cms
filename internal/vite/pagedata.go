package vite

import (
	"html/template"
)

type PageData struct {
	IsDev          bool
	ViteURL        string
	Metadata       template.HTML
	StyleSheets    template.HTML
	Modules        template.HTML
	PreloadModules template.HTML
	Scripts        template.HTML
	Content        template.HTML

  Header template.HTML
}
