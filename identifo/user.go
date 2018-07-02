package identifo

//UserID is datatype to store unique user identifier
type UserID string

//User is general model to represent the user data
type User struct {
	ID      UserID                 `json:"id,omitempty"`
	Name    string                 `json:"name,omitempty"`
	Phone   string                 `json:"phone,omitempty"`
	Email   string                 `json:"email,omitempty"`
	Profile map[string]interface{} `json:"profile,omitempty"`
}
