package main

import (
	"fmt"
	"log"
	"net/http"

	ihttp "github.com/madappgang/identifo/http"
	"github.com/madappgang/identifo/jwt"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/mongo"
)

func initServices() (model.AppStorage, model.UserStorage, model.TokenStorage, model.TokenService) {
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

	if _, err = appStorage.AppByID("59fd884d8f6b180001f5b4e2"); err != nil {
		fmt.Printf("Creating data because got error trying to get app: %+v\n", err)
		createData(db, userStorage.(*mongo.UserStorage), appStorage)
	}
	return appStorage, userStorage, tokenStorage, tokenService
}

func staticPages() ihttp.StaticPages {
	return ihttp.StaticPages{
		Login:          "../../static/login.html",
		Registration:   "../../static/registration.html",
		ForgotPassword: "../../static/forgot-password.html",
		ResetPassword:  "../../static/reset-password.html",
	}
}

func staticFiles() ihttp.StaticFiles {
	return ihttp.StaticFiles{
		StylesDirectory:  "../../static/css",
		ScriptsDirectory: "../../static/js",
	}
}

func initRouter() model.Router {
	appStorage, userStorage, tokenStorage, tokenService := initServices()

	sp := staticPages()
	sf := staticFiles()

	router, err := ihttp.NewRouter(nil, appStorage, userStorage, tokenStorage, tokenService, ihttp.ServeStaticPages(sp), ihttp.ServeStaticFiles(sf))

	if err != nil {
		log.Fatal(err)
	}

	return router
}

func main() {
	fmt.Println("mongodb server started")
	r := initRouter()

	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}

func createData(db *mongo.DB, us *mongo.UserStorage, as model.AppStorage) {
	u1d := []byte(`{"name":"test@madappgang.com","active":true}`)
	u1, err := mongo.UserFromJSON(u1d)
	if err != nil {
		log.Fatal(err)
	}
	if _, err = us.AddNewUser(u1, "secret"); err != nil {
		log.Fatal(err)
	}

	u1d = []byte(`{"name":"User2","active":false}`)
	u1, err = mongo.UserFromJSON(u1d)
	if err != nil {
		log.Fatal(err)
	}
	if _, err = us.AddNewUser(u1, "other_password"); err != nil {
		log.Fatal(err)
	}

	ad := []byte(`{
		"id":"59fd884d8f6b180001f5b4e2",
		"secret":"secret",
		"name":"iOS App",
		"active":true, 
		"description":"Amazing ios app", 
		"scopes":["smartrun"],
		"offline":true,
		"redirect_url":"myapp://loginhook",
		"refresh_token_lifespan":9000000,
		"token_lifespan":9000
	}`)
	app, err := mongo.AppDataFromJSON(ad)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("app data: %+v", app)
	if _, err = as.AddNewApp(app); err != nil {
		log.Fatal(err)
	}
}
