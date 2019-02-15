package server

import (
	"github.com/labstack/echo"
	"html/template"
	"io"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// NewTemplateRenderer return new
func NewTemplateRenderer(templates *template.Template) *TemplateRenderer {
	return &TemplateRenderer{
		templates: templates,
	}
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
