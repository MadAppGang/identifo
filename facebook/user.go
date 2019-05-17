package facebook

// User is user profile data on Facebook platform: https://developers.facebook.com/docs/graph-api/reference/v2.6/user.
type User struct {
	ID         string `json:"id,omitempty"` //The id of this person's user account. This ID is unique to each app and cannot be used across different apps.
	Email      string `json:"email,omitempty"`
	Name       string `json:"name,omitempty"`
	ProfilePic string `json:"profile_pic,omitempty"`
}
