package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

<<<<<<< HEAD
	"github.com/joho/godotenv"
	"github.com/madappgang/identifo/dynamodb"
	ihttp "github.com/madappgang/identifo/http"
	"github.com/madappgang/identifo/jwt"
	"github.com/madappgang/identifo/mailgun"
=======
>>>>>>> 110cf49475c488a82de8b113bee667d971b4b81e
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server/dynamodb"
)

<<<<<<< HEAD
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

func initDB() model.Router {
	db, err := dynamodb.NewDB("http://localhost:8000", "")
	if err != nil {
		log.Fatal(err)
	}
	appStorage, err := dynamodb.NewAppStorage(db)
	if err != nil {
		log.Fatal(err)
	}

	userStorage, err := dynamodb.NewUserStorage(db)
	if err != nil {
		log.Fatal(err)
	}

	tokenStorage, err := dynamodb.NewTokenStorage(db)
	if err != nil {
		log.Fatal(err)
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	domain := os.Getenv("MAILGUN_DOMAIN")
	privateKey := os.Getenv("MAILGUN_PRIVATE_KEY")
	publicKey := os.Getenv("MAILGUN_PUBLIC_KEY")
	emailService := mailgun.NewEmailService(domain, privateKey, publicKey, "sender@identifo.com")

	tokenService, err := jwt.NewTokenService(
		"../../jwt/private.pem",
		"../../jwt/public.pem",
		"http://localhost:8080",
		model.TokenServiceAlgorithmAuto,
		tokenStorage,
		appStorage,
		userStorage,
	)

	sp := staticPages()
	sf := staticFiles()

	r, err := ihttp.NewRouter(nil, appStorage, userStorage, tokenStorage, tokenService, emailService, ihttp.ServeStaticPages(sp), ihttp.ServeStaticFiles(sf))
=======
func initDB() model.Server {
	settings := dynamodb.DefaultSettings
	settings.StaticFolderPath = "../.."
	settings.PEMFolderPath = "../../jwt"
	settings.Issuer = "http://localhost:8080"
>>>>>>> 110cf49475c488a82de8b113bee667d971b4b81e

	server, err := dynamodb.NewServer(settings)
	if err != nil {
		log.Fatal(err)
	}

	_, err = server.AppStorage().AppByID("59fd884d8f6b180001f5b4e2")
	if err != nil {
		server.ImportApps("../import/apps.json")
		server.ImportUsers("../import/users.json")
	}
	return server
}

func main() {
	r := initDB()
	fmt.Println("DynamoDB server started")
	log.Fatal(http.ListenAndServe(":8080", r.Router()))
}
