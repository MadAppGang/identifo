package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/madappgang/identifo/dynamodb"
	"github.com/madappgang/identifo/jwt"
	"github.com/madappgang/identifo/mailgun"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/web"
	"github.com/madappgang/identifo/web/html"
)

func staticPath() html.StaticFilesPath {
	return html.StaticFilesPath{
		StylesPath:  "../../static/css",
		ScriptsPath: "../../static/js",
		PagesPath:   "../../static",
		ImagesPath:  "../../static/img",
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

	_ = godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("Error loading .env file: %v", err)
	// }

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

	routerSettings := web.RouterSetting{
		AppStorage:   appStorage,
		UserStorage:  userStorage,
		TokenStorage: tokenStorage,
		TokenService: tokenService,
		EmailService: emailService,
		WebRouterSettings: []func(*html.Router) error{
			html.StaticPathOptions(staticPath()),
		},
	}
	r, err := web.NewRouter(routerSettings)

	if err != nil {
		log.Fatal(err)
	}

	_, err = appStorage.AppByID("59fd884d8f6b180001f5b4e2")
	if err != nil {
		fmt.Printf("Creating data because got error trying to get app: %+v\n", err)
		createData(db, userStorage.(*dynamodb.UserStorage), appStorage)
	}
	return r
}

func main() {
	r := initDB()
	fmt.Println("dynamoDB server started")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func createData(db *dynamodb.DB, us *dynamodb.UserStorage, as model.AppStorage) {
	u1d := []byte(`{"username":"test@madappgang.com","active":true}`)
	u1, _ := dynamodb.UserFromJSON(u1d)
	us.AddNewUser(u1, "secret")

	u1d = []byte(`{"username":"User2","active":false}`)
	u1, _ = dynamodb.UserFromJSON(u1d)
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
	app, err := dynamodb.AppDataFromJSON(ad)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("app data: %+v", app)
	_, err = as.AddNewApp(app)
	if err != nil {
		log.Fatal(err)
	}
}
