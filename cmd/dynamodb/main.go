package main

import (
	"fmt"
	"log"
	"net/http"

	ddb "github.com/madappgang/identifo/dynamodb"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server/dynamodb"
)

func initDB() model.Server {
	settings := dynamodb.DefaultSettings
	settings.StaticFolderPath = "../.."
	settings.PEMFolderPath = "../../jwt"
	settings.Issuer = "http://localhost:8080"

	server, err := dynamodb.NewServer(settings)
	if err != nil {
		log.Fatal(err)
	}

	_, err = server.AppStorage().AppByID("59fd884d8f6b180001f5b4e2")
	if err != nil {
		fmt.Printf("Creating data because got error trying to get app: %+v\n", err)
		createData(server.UserStorage().(*ddb.UserStorage), server.AppStorage())
	}
	return server
}

func main() {
	r := initDB()
	fmt.Println("dynamoDB server started")
	log.Fatal(http.ListenAndServe(":8080", r.Router()))
}

func createData(us *ddb.UserStorage, as model.AppStorage) {
	u1d := []byte(`{"username":"test@madappgang.com","active":true}`)
	u1, _ := ddb.UserFromJSON(u1d)
	us.AddNewUser(u1, "secret")

	u1d = []byte(`{"username":"User2","active":false}`)
	u1, _ = ddb.UserFromJSON(u1d)
	us.AddNewUser(u1, "other_password")

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
	app, err := ddb.AppDataFromJSON(ad)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("app data: %+v", app)
	_, err = as.AddNewApp(app)
	if err != nil {
		log.Fatal(err)
	}
}
