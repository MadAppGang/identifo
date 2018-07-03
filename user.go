package identifo

//UserID is datatype to store unique user identifier
type UserID string

//User is general model to represent the user data
type User struct {
	ID      UserID                 `json:"id,omitempty"`
	Profile map[string]interface{} `json:"profile,omitempty"`
}
