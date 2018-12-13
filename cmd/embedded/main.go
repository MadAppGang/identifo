package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/madappgang/identifo/boltdb"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server/embedded"
)

func initDB() model.Server {
	settings := embedded.DefaultSettings
	settings.StaticFolderPath = "../.."
	settings.PEMFolderPath = "../../jwt"
	settings.Issuer = "http://localhost:8080"

	server, err := embedded.NewServer(settings)
	if err != nil {
		log.Fatal(err)
	}

	_, err = server.AppStorage().AppByID("59fd884d8f6b180001f5b4e2")
	if err != nil {
		fmt.Printf("Creating data because got error trying to get app: %+v\n", err)
		createData(server.UserStorage().(*boltdb.UserStorage), server.AppStorage())
	}
	return server
}
func main() {
	r := initDB()
	fmt.Println("Embedded server started")
	log.Fatal(http.ListenAndServe(":8080", r.Router()))
}

func createData(us *boltdb.UserStorage, as model.AppStorage) {
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
