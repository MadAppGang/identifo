package main

import (
	"log"
	"net/http"
	"os"

<<<<<<< HEAD
	"github.com/joho/godotenv"
	ihttp "github.com/madappgang/identifo/http"
	"github.com/madappgang/identifo/jwt"
	"github.com/madappgang/identifo/mailgun"
	"github.com/madappgang/identifo/mem"
	"github.com/madappgang/identifo/model"
=======
	"github.com/madappgang/identifo/server/fake"
>>>>>>> 110cf49475c488a82de8b113bee667d971b4b81e
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

<<<<<<< HEAD
	appStorage := mem.NewAppStorage()
	userStorage := mem.NewUserStorage()
	tokenStorage := mem.NewTokenStorage()

	app := mem.MakeAppData("59fd884d8f6b180001f5b4e2", "secret", true, "Test app", []string{"offline", "smartrun"}, true, "", 0, 0)
	if _, err := appStorage.AddNewApp(app); err != nil {
		panic(err)
	}

	tokenService, _ := jwt.NewTokenService(
		"../../jwt/private.pem",
		"../../jwt/public.pem",
		"identifo.madappgang.com",
		model.TokenServiceAlgorithmAuto,
		tokenStorage,
		appStorage,
		userStorage,
	)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	domain := os.Getenv("MAILGUN_DOMAIN")
	privateKey := os.Getenv("MAILGUN_PRIVATE_KEY")
	publicKey := os.Getenv("MAILGUN_PUBLIC_KEY")
	emailService := mailgun.NewEmailService(domain, privateKey, publicKey, "sender@identifo.com")

	sp := staticPages()
	sf := staticFiles()

	r, err := ihttp.NewRouter(nil, appStorage, userStorage, tokenStorage, tokenService, emailService, ihttp.ServeStaticPages(sp), ihttp.ServeStaticFiles(sf))
=======
	settings := fake.DefaultSettings
	settings.StaticFolderPath = "../.."
	settings.PEMFolderPath = "../../jwt"
	settings.Issuer = "http://localhost:8080"
>>>>>>> 110cf49475c488a82de8b113bee667d971b4b81e

	server, err := fake.NewServer(settings)
	if err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":8080", server.Router()); err != nil {
		panic(err)
	}
}
