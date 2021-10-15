package model

import "html/template"

type TemplateStorage interface {
	Template(name string) (*template.Template, error)
}
