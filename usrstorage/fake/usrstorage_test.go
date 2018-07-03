package fake_test

import (
	"testing"

	"github.com/MadAppGang/identifo"
	"github.com/MadAppGang/identifo/usrstorage/fake"
)

func TestUsrstorage_CreateClient(t *testing.T) {
	c := fake.NewClient(nil)
	if c == nil {
		t.Fatal("the client is empty")
	}
}

func TestUsrstorage_CreateSession(t *testing.T) {
	c := fake.NewClient(nil)
	if c == nil {
		t.Fatal("the client is empty")
	}

	s, err := c.Connect()
	if err != nil {
		t.Fatal(err)
	}
	if s == nil {
		t.Fatal("session is empty")
	}

	st := s.Storage()
	if st == nil {
		t.Fatal("storage is empty")
	}
}

func TestUsrstorage_GetUser(t *testing.T) {
	users := []identifo.User{
		identifo.User{ID: "id123", Profile: map[string]interface{}{fake.PasswordKey: "password"}},
	}

	c := fake.NewClient(users)
	if c == nil {
		t.Fatal("the client is empty")
	}

	s, err := c.Connect()
	if err != nil {
		t.Fatal(err)
	}
	if s == nil {
		t.Fatal("session is empty")
	}

	u, err := s.Storage().FindUser("id123", "password")
	if err != nil {
		t.Fatal(err)
	}
	if u == nil {
		t.Fatal("empty user found")
	}
	if u.ID != "id123" {
		t.Fatal("wrong user")
	}

}

func TestUsrstorage_GetUserByEmail(t *testing.T) {
	users := []identifo.User{
		identifo.User{ID: "id123", Profile: map[string]interface{}{fake.PasswordKey: "password", "email": "mail@mail.com"}},
	}

	c := fake.NewClient(users)
	if c == nil {
		t.Fatal("the client is empty")
	}

	s, err := c.Connect()
	if err != nil {
		t.Fatal(err)
	}
	if s == nil {
		t.Fatal("session is empty")
	}

	u, err := s.Storage().FindUserWithKey("email", "mail@mail.com", "password")
	if err != nil {
		t.Fatal(err)
	}
	if u == nil {
		t.Fatal("empty user found")
	}
	if u.ID != "id123" || u.Profile["email"] != "mail@mail.com" {
		t.Fatal("wrong user")
	}
}

func TestUsrstorage_CreateUser(t *testing.T) {
	users := []identifo.User{
		identifo.User{ID: "id123", Profile: map[string]interface{}{fake.PasswordKey: "password", "email": "mail@mail.com"}},
	}
	c := fake.NewClient(users)
	if c == nil {
		t.Fatal("the client is empty")
	}

	s, err := c.Connect()
	if err != nil {
		t.Fatal(err)
	}
	if s == nil {
		t.Fatal("session is empty")
	}

	newUser := identifo.User{Profile: map[string]interface{}{"email": "mail2@mail.com"}}
	createdUser, err := s.Storage().CreateUser(newUser, "password2")
	if err != nil {
		t.Fatal(err)
	}
	if createdUser == nil {
		t.Fatal("no user created")
	}
	if createdUser.ID == "" || createdUser.Profile[fake.PasswordKey] == nil || createdUser.Profile["email"] != newUser.Profile["email"] {
		t.Fatal("new user has wrong metadata")
	}

	u, err := s.Storage().FindUserWithKey("email", "mail2@mail.com", "password2")
	if err != nil {
		t.Fatal(err)
	}
	if u == nil {
		t.Fatal("empty user found")
	}
	if u.ID != createdUser.ID || u.Profile["email"] != "mail2@mail.com" {
		t.Fatal("wrong user")
	}
}

func TestUsrstorage_GetUser_ErrorNotFound(t *testing.T) {
	users := []identifo.User{
		identifo.User{ID: "id123", Profile: map[string]interface{}{fake.PasswordKey: "password"}},
	}

	c := fake.NewClient(users)
	if c == nil {
		t.Fatal("the client is empty")
	}

	s, err := c.Connect()
	if err != nil {
		t.Fatal(err)
	}
	if s == nil {
		t.Fatal("session is empty")
	}

	u, err := s.Storage().FindUser("id123", "wrong_password")
	if u != nil {
		t.Fatal("user should be empty")
	}
	if err == nil {
		t.Fatal("error should be here")
	}
	if err != fake.ErrorUserNotFound {
		t.Fatal(err)
	}

	u, err = s.Storage().FindUser("wrong_id", "password")
	if u != nil {
		t.Fatal("user should be empty")
	}
	if err == nil {
		t.Fatal("error should be here")
	}
	if err != fake.ErrorUserNotFound {
		t.Fatal(err)
	}

}
