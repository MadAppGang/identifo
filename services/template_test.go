package services_test

import (
	"bytes"
	"testing"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/services"
	"github.com/madappgang/identifo/storage"
)

func TestFSTemplate_Template(t *testing.T) {
	settings := model.FileStorageSettings{
		Type: model.FileStorageTypeLocal,
		Local: model.FileStorageLocal{
			FolderPath: "../test/artifacts/templates",
		},
	}

	fss, err := storage.NewFS(settings)
	if err != nil {
		t.Fatalf("error creating local fs with error: %v", err)
	}

	templateStorage := services.NewTemplate(fss)
	tml, err := templateStorage.Template("mail1.template")
	if err != nil {
		t.Fatalf("error creating email 1 template: %v", err)
	}

	buf := bytes.Buffer{}
	data := struct {
		Title string
	}{
		Title: "1",
	}
	tml.Execute(&buf, data)
	expectedRender := "I am email template 1"
	if buf.String() != expectedRender {
		t.Fatalf("wrong template rendering, got: %s, expected: %s", buf.String(), expectedRender)
	}

	// check template from cache, reading second time should return cached template
	tml2, err := templateStorage.Template("mail1.template")
	if err != nil {
		t.Fatalf("error creating email 1 template: %v", err)
	}
	buf2 := bytes.Buffer{}
	tml2.Execute(&buf2, data)
	if buf.String() != expectedRender {
		t.Fatalf("wrong template rendering, got: %s, expected: %s", buf.String(), expectedRender)
	}
}
