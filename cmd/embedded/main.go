package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/madappgang/identifo/boltdb"
	ihttp "github.com/madappgang/identifo/http"
	"github.com/madappgang/identifo/jwt"
	"github.com/madappgang/identifo/model"
)

func initServices() (model.AppStorage, model.UserStorage, model.TokenStorage, model.TokenService) {
	db, err := boltdb.InitDB("db.db")
	if err != nil {
		log.Fatal(err)
	}
	appStorage, _ := boltdb.NewAppStorage(db)
	userStorage, _ := boltdb.NewUserStorage(db)
	tokenStorage, _ := boltdb.NewTokenStorage(db)

	tokenService, _ := jwt.NewTokenService(
		"../../jwt/private.pem",
		"../../jwt/public.pem",
		"identifo.madappgang.com",
		model.TokenServiceAlgorithmAuto,
		tokenStorage,
		appStorage,
		userStorage,
	)

	_, err = appStorage.AppByID("59fd884d8f6b180001f5b4e2")

	if err != nil {
		fmt.Printf("Creating data because got error trying to get app: %+v\n", err)
		createData(db, userStorage.(*boltdb.UserStorage), appStorage)
	}

	return appStorage, userStorage, tokenStorage, tokenService
}

func initRouter() model.Router {
	appStorage, userStorage, tokenStorage, tokenService := initServices()

	router, err := ihttp.NewRouter(nil, appStorage, userStorage, tokenStorage, tokenService)

	if err != nil {
		log.Fatal(err)
	}

	return router
}

func main() {
	fmt.Println("Embedded server started")
	r := initRouter()

	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}

func createData(db *bolt.DB, us *boltdb.UserStorage, as model.AppStorage) {
	u1d := []byte(`{"id":"12345","name":"test@madappgang.com","active":true}`)
	u1, err := boltdb.UserFromJSON(u1d)
	if err != nil {
		log.Fatal(err)
	}
	if _, err = us.AddNewUser(u1, "secret"); err != nil {
		log.Fatal(err)
	}

	u1d = []byte(`{"id":"12346","name":"User2","active":false}`)
	u1, err = boltdb.UserFromJSON(u1d)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := us.AddNewUser(u1, "other_password"); err != nil {
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
		"reset_password_token_lifespan": 7200,
		"token_lifespan":9000
	}`)
	app, err := boltdb.AppDataFromJSON(ad)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("app data: %+v", app)
	if _, err = as.AddNewApp(app); err != nil {
		log.Fatal(err)
	}
}
