package mail_test

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/services/mail"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

type MessageProvider interface {
	Messages() []map[string]string
}

func TestEmailService_SendMessage(t *testing.T) {
	sts := model.EmailServiceSettings{
		Type: model.EmailServiceMock,
	}

	afs := createFS()
	service, err := mail.NewService(sts, afero.NewIOFS(afs), time.Second, "templates")
	require.NoError(t, err)
	service.Start()
	err = service.SendUserEmail(
		model.EmailTemplateTypeVerifyEmail,
		"",
		model.User{
			Email:    "mail@mail.com",
			Locale:   "",
			Username: "Jack Daniels",
		},
		map[string]string{
			"Message": "Hello World",
		},
	)
	require.NoError(t, err)

	var tr MessageProvider = service.Transport().(MessageProvider)
	assert.Contains(t, tr.Messages()[0]["body"], "I AM INITIAL TEMPLATE")
	assert.NotContains(t, tr.Messages()[0]["body"], "I AM UPDATED TEMPLATE")
	assert.Contains(t, tr.Messages()[0]["body"], "Hello World")
	assert.Contains(t, tr.Messages()[0]["subject"], "Jack Daniels")

	// let's change the template
	afs.Remove(filepath.Join("templates", model.EmailTemplateTypeVerifyEmail.FileName()))
	body := `---
I am updated template subject {{.User.Username}}
---
<html>
	<body>
		<h1>I AM UPDATED TEMPLATE</h1>
	</body>
</html>
---
`
	afero.WriteFile(afs, filepath.Join("templates", model.EmailTemplateTypeVerifyEmail.FileName()), []byte(body), 0o644)

	time.Sleep(time.Second * 2)
	err = service.SendUserEmail(
		model.EmailTemplateTypeVerifyEmail,
		"",
		model.User{
			Email:  "mail@mail.com",
			Locale: "",
		},
		model.EmailData{},
	)
	require.NoError(t, err)
	assert.NotContains(t, tr.Messages()[1]["body"], "I AM INITIAL TEMPLATE")
	assert.Contains(t, tr.Messages()[1]["body"], "I AM UPDATED TEMPLATE")
}

func TestEmailService_LocalizedEmail(t *testing.T) {
	sts := model.EmailServiceSettings{
		Type: model.EmailServiceMock,
	}

	afs := createFS()
	service, err := mail.NewService(sts, afero.NewIOFS(afs), time.Second, "templates")
	require.NoError(t, err)
	service.Start()
	err = service.SendUserEmail(
		model.EmailTemplateTypeVerifyEmail,
		"",
		model.User{
			Email:    "mail@mail.com",
			Locale:   language.Ukrainian.String(),
			Username: "Jack Daniels",
		},
		map[string]string{
			"Message": "Доброго ранку!",
		},
	)
	require.NoError(t, err)
	fmt.Println(language.Ukrainian.String())
	var tr MessageProvider = service.Transport().(MessageProvider)
	assert.Contains(t, tr.Messages()[0]["body"], "Я Українська версія листа")
	assert.Contains(t, tr.Messages()[0]["body"], "Доброго ранку!")
	assert.Contains(t, tr.Messages()[0]["subject"], "Jack Daniels")
	assert.Contains(t, tr.Messages()[0]["subject"], "Тема спілкування")

	err = service.SendUserEmail(
		model.EmailTemplateTypeVerifyEmail,
		"",
		model.User{
			Email:  "mail@mail.com",
			Locale: language.Albanian.String(),
		},
		map[string]string{
			"Message": "Доброго ранку!",
		},
	)
	require.NoError(t, err)
	assert.Contains(t, tr.Messages()[1]["body"], "I AM INITIAL TEMPLATE")
	assert.Contains(t, tr.Messages()[1]["subject"], "I am template subject")
}

func createFS() afero.Fs {
	body := `---
I am template subject {{.User.Username}}
---
<html>
	<body>
		<h1>I AM INITIAL TEMPLATE</h1>
		<h1>{{.Data.Message}}</h1>
	</body>
</html>
---
`

	ukbody := `---
Тема спілкування {{.User.Username}}
---
<html>
	<body>
		<h1>Я Українська версія листа</h1>
		<h1>{{.Data.Message}}</h1>
	</body>
</html>
---
`
	tplFS := afero.NewMemMapFs()
	// create test files and directories
	tplFS.MkdirAll("templates", 0o755)
	for _, tpl := range model.AllEmailTemplatesFileNames() {
		afero.WriteFile(tplFS, filepath.Join("templates", tpl), []byte(body), 0o644)
		afero.WriteFile(tplFS, filepath.Join("templates", tpl+"_"+language.Ukrainian.String()), []byte(ukbody), 0o644)
	}
	return tplFS
}
