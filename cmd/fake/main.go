package main

import (
	"log"
	"net/http"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server/fake"
)

func staticPages() ihttp.StaticPages {
	return ihttp.StaticPages{
		Login:                 "../../static/login.html",
		Registration:          "../../static/registration.html",
		ForgotPassword:        "../../static/forgot-password.html",
		ResetPassword:         "../../static/reset-password.html",
		ForgotPasswordSuccess: "../../static/forgot-password-success.html",
		TokenError:            "../../static/token-error.html",
		ResetSuccess:          "../../static/reset-success.html",
	}
}

func staticFiles() ihttp.StaticFiles {
	return ihttp.StaticFiles{
		StylesDirectory:  "../../static/css",
		ScriptsDirectory: "../../static/js",
	}
}

//this server works only with memory storages and generated data
//should be used for test and CI environments only
func main() {

	settings := fake.DefaultSettings
	settings.StaticFolderPath = "../../static"
	settings.PEMFolderPath = "../../jwt"
	settings.Issuer = "http://localhost:8080"
	settings.MailService = model.MailServiceAWS
	settings.EmailTemplatesPath = "../../email_templates"

	server, err := fake.NewServer(settings)
	if err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":8080", server.Router()); err != nil {
		panic(err)
	}
}
