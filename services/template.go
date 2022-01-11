package services

import (
	"html/template"
	"io/fs"

	"github.com/madappgang/identifo/v2/model"
)

func NewTemplate(fss fs.FS) model.TemplateStorage {
	c := make(map[string]template.Template)
	t := FSTemplate{fss: fss, cache: c}
	return &t
}

// Template is default template implementation, that uses fs.FS as storage backend
type FSTemplate struct {
	fss   fs.FS
	cache map[string]template.Template
}

func (t *FSTemplate) Template(name string) (*template.Template, error) {
	tt, ok := t.cache[name]

	if ok == false {
		tmpl, err := template.ParseFS(t.fss, name)
		if err != nil {
			return nil, err
		}
		tt = *tmpl
		t.cache[name] = tt
	}
	return &tt, nil
}
