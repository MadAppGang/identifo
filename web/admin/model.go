package admin

import "github.com/madappgang/identifo/v2/model"


// reset email data????
type resetEmailData struct {
	User  model.User
	Token string
	URL   string
	Host  string
}
