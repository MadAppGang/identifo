package main

import (
	"fmt"
	"net/http"

	ihttp "github.com/madappgang/identifo/http"
	"github.com/madappgang/identifo/jwt"
	"github.com/madappgang/identifo/mem"
)

//this server works only with memory storages and generated data
//should be used for test and CI environments only
func main() {
	fmt.Println("Fake server started")

	appStorage := mem.NewAppStorage()
	userStorage := mem.NewUserStorage()
	tokenStorage := mem.NewTokenStorage()

	app := mem.MakeAppData("12345", "secret", true, "Test app", []string{"offline", "smartrun"}, true, "", 0, 0)
	appStorage.AddNewApp(app)

	tokenService, _ := jwt.NewTokenService(
		"../../jwt/private.pem",
		"../../jwt/public.pem",
		"identifo.madappgang.com",
		tokenStorage,
		appStorage,
		userStorage,
	)
	r := ihttp.NewRouter(nil, appStorage, userStorage, tokenStorage, tokenService)
	http.ListenAndServe(":8080", r)
}
