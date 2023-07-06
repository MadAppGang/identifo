package mail_test

type MessageProvider interface {
	Messages() []string
}

// func TestEmailService_SendMessage(t *testing.T) {
// 	sts := model.EmailServiceSettings{
// 		Type: model.EmailServiceMock,
// 	}

// 	afs := createFS()
// 	service, err := mail.NewService(sts, afero.NewIOFS(afs), time.Second, "templates")
// 	require.NoError(t, err)
// 	service.Start()
// 	service.SendTemplateEmail(
// 		model.EmailTemplateTypeVerifyEmail,
// 		"",
// 		"SUBJECT1",
// 		"mail@mail.com",
// 		model.EmailData{},
// 	)

// 	var tr MessageProvider = service.Transport().(MessageProvider)
// 	assert.Contains(t, tr.Messages()[0], "I AM INITIAL TEMPLATE")
// 	assert.NotContains(t, tr.Messages()[0], "I AM UPDATED TEMPLATE")

// 	// let's change the template
// 	afs.Remove(filepath.Join("templates", model.EmailTemplateTypeVerifyEmail.FileName()))
// 	afero.WriteFile(afs, filepath.Join("templates", model.EmailTemplateTypeVerifyEmail.FileName()), []byte("I AM UPDATED TEMPLATE"), 0o644)
// 	time.Sleep(time.Second * 2)
// 	service.SendTemplateEmail(
// 		model.EmailTemplateTypeVerifyEmail,
// 		"",
// 		"SUBJECT1",
// 		"mail@mail.com",
// 		model.EmailData{},
// 	)
// 	fmt.Println(tr.Messages()[0])
// 	fmt.Println(tr.Messages()[1])
// 	assert.NotContains(t, tr.Messages()[1], "I AM INITIAL TEMPLATE")
// 	assert.Contains(t, tr.Messages()[1], "I AM UPDATED TEMPLATE")
// }

// func createFS() afero.Fs {
// 	tplFS := afero.NewMemMapFs()
// 	// create test files and directories
// 	tplFS.MkdirAll("templates", 0o755)
// 	for _, tpl := range model.AllEmailTemplatesFileNames() {
// 		afero.WriteFile(tplFS, filepath.Join("templates", tpl), []byte("I AM INITIAL TEMPLATE"), 0o644)
// 	}
// 	return tplFS
// }
