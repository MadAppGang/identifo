package model

import (
	"io/fs"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Server holds together all dependencies.
type Server interface {
	Router() Router
	UpdateCORS()
	Storages() ServerStorageCollection
	Services() ServerServices
	Settings() ServerSettings
	Errors() []error
	Close()
}

// Router handles all incoming http requests.
type Router interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// ServerStorageCollection holds the full collections of server storage components
type ServerStorageCollection struct {
	App           AppStorage
	User          UserStorage
	Token         TokenStorage
	Blocklist     TokenBlacklist
	Invite        InviteStorage
	Verification  VerificationCodeStorage
	Config        ConfigurationStorage
	Session       SessionStorage
	Key           KeyStorage
	ManagementKey ManagementKeysStorage
	LoginAppFS    fs.FS
	AdminPanelFS  fs.FS
}

type ServerServices struct {
	SMS     SMSService
	Email   EmailService
	Token   TokenService
	Session SessionService
}

// PasswordHash creates hash with salt for password.
func PasswordHash(pwd string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash)
}

// RandomPassword creates random password
func RandomPassword(length int) string {
	rand.Seed(time.Now().UnixNano())
	return randSeq(length)
}

var rndPassLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890?!@#")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = rndPassLetters[rand.Intn(len(rndPassLetters))]
	}
	return string(b)
}
