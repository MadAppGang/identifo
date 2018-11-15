package main

import (
	"fmt"
	"log"
	"net/http"

	ihttp "github.com/madappgang/identifo/http"
	"github.com/madappgang/identifo/jwt"
	"github.com/madappgang/identifo/mem"
	"github.com/madappgang/identifo/model"
)

//this server works only with memory storages and generated data
//should be used for test and CI environments only
func main() {
	fmt.Println("Fake server started")

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

	r, err := ihttp.NewRouter(nil, appStorage, userStorage, tokenStorage, tokenService)

	if err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}
