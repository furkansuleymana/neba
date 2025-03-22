package ui

import (
	"embed"
	"io/fs"
)

//go:embed all:templates
var dir embed.FS
var TemplatesDirFS, _ = fs.Sub(dir, "templates")
