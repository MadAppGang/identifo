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

	app := mem.MakeAppData("59fd884d8f6b180001f5b4e2", "secret", true, "Test app", []string{"offline", "smartrun"}, true, "", 0, 0)
	appStorage.AddNewApp(app)

	tokenService, _ := jwt.NewTokenService(
		"../../jwt/private.pem",
		"../../jwt/public.pem",
		"identifo.madappgang.com",
		tokenStorage,
		appStorage,
		userStorage,
	)

	var settings ihttp.Settings
	r := ihttp.NewRouter(nil, appStorage, userStorage, tokenStorage, tokenService, settings)
	http.ListenAndServe(":8080", r)
}
