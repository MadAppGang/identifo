package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/madappgang/identifo/dynamodb"
	ihttp "github.com/madappgang/identifo/http"
	"github.com/madappgang/identifo/jwt"
	"github.com/madappgang/identifo/model"
)

func initDB() model.Router {
	db, err := dynamodb.NewDB("http://localhost:8000", "")
	if err != nil {
		log.Fatal(err)
	}
	appStorage, _ := dynamodb.NewAppStorage(db)
	userStorage, _ := dynamodb.NewUserStorage(db)
	tokenStorage, _ := dynamodb.NewTokenStorage(db)

	tokenService, _ := jwt.NewTokenService(
		"../../jwt/private.pem",
		"../../jwt/public.pem",
		"identifo.madappgang.com",
		tokenStorage,
		appStorage,
		userStorage,
	)
	r, err := ihttp.NewRouter(nil, appStorage, userStorage, tokenStorage, tokenService)

	if err != nil {
		log.Fata(err)
	}

	_, err = appStorage.AppByID("59fd884d8f6b180001f5b4e2")
	if err != nil {
		fmt.Printf("Creating data because got error trying to get app: %+v\n", err)
		createData(db, userStorage.(*dynamodb.UserStorage), appStorage)
	}
	return r
}

func main() {
	fmt.Println("dynamoDB server started")
	r := initDB()

	log.Fatal(http.ListenAndServe(":8080", r))
}

func createData(db *dynamodb.DB, us *dynamodb.UserStorage, as model.AppStorage) {
	u1d := []byte(`{"username":"test@madappgang.com","active":true}`)
	u1, err := dynamodb.UserFromJSON(u1d)
	if err != nil {
		log.Fatal(err)
	}
	if _, err = us.AddNewUser(u1, "secret"); err != nil {
		log.Fatal(err)
	}

	u1d = []byte(`{"username":"User2","active":false}`)
	u1, err = dynamodb.UserFromJSON(u1d)
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
		"token_lifespan":9000
	}`)
	app, err := dynamodb.AppDataFromJSON(ad)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("app data: %+v", app)
	if _, err = as.AddNewApp(app); err != nil {
		log.Fatal(err)
	}
}
