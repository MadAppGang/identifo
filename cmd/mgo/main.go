package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

<<<<<<< HEAD
	"github.com/joho/godotenv"
	ihttp "github.com/madappgang/identifo/http"
	"github.com/madappgang/identifo/jwt"
	"github.com/madappgang/identifo/mailgun"
=======
>>>>>>> 110cf49475c488a82de8b113bee667d971b4b81e
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server/mgo"
)

<<<<<<< HEAD
func initServices() (model.AppStorage, model.UserStorage, model.TokenStorage, model.TokenService, model.EmailService) {
	db, err := mongo.NewDB("localhost:27017", "identifo")
	if err != nil {
		log.Fatal(err)
	}
	appStorage, _ := mongo.NewAppStorage(db)
	userStorage, _ := mongo.NewUserStorage(db)
	tokenStorage, _ := mongo.NewTokenStorage(db)

	tokenService, _ := jwt.NewTokenService(
		"../../jwt/private.pem",
		"../../jwt/public.pem",
		"identifo.madappgang.com",
		model.TokenServiceAlgorithmAuto,
		tokenStorage,
		appStorage,
		userStorage,
	)

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	domain := os.Getenv("MAILGUN_DOMAIN")
	privateKey := os.Getenv("MAILGUN_PRIVATE_KEY")
	publicKey := os.Getenv("MAILGUN_PUBLIC_KEY")
	emailService := mailgun.NewEmailService(domain, privateKey, publicKey, "sender@identifo.com")

	if _, err = appStorage.AppByID("59fd884d8f6b180001f5b4e2"); err != nil {
		fmt.Printf("Creating data because got error trying to get app: %+v\n", err)
		createData(db, userStorage.(*mongo.UserStorage), appStorage)
	}
	return appStorage, userStorage, tokenStorage, tokenService, emailService
}

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

func initRouter() model.Router {
	appStorage, userStorage, tokenStorage, tokenService, emailService := initServices()

	sp := staticPages()
	sf := staticFiles()

	router, err := ihttp.NewRouter(nil, appStorage, userStorage, tokenStorage, tokenService, emailService, ihttp.ServeStaticPages(sp), ihttp.ServeStaticFiles(sf))
=======
func server() model.Server {
	settings := mgo.DefaultSettings
	settings.StaticFolderPath = "../.."
	settings.PEMFolderPath = "../../jwt"
	settings.Issuer = "http://localhost:8080"
	settings.DBEndpoint = "localhost:27017"
	settings.DBName = "identifo"

	server, err := mgo.NewServer(settings)
	if err != nil {
		log.Fatal(err)
	}
>>>>>>> 110cf49475c488a82de8b113bee667d971b4b81e

	_, err = server.AppStorage().AppByID("59fd884d8f6b180001f5b4e2")
	if err != nil {
		server.ImportApps("../import/apps.json")
		server.ImportUsers("../import/users.json")
	}
	return server
}

func main() {
	fmt.Println("mongodb server started")
	r := server()

	if err := http.ListenAndServe(":8080", r.Router()); err != nil {
		panic(err)
	}
}
