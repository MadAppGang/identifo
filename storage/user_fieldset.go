package storage

type UserFieldset string

const (
	UserFieldsetBasic UserFieldset = "basic"
	UserFieldsetAll   UserFieldset = "all"
	// TODO: Add fieldset cases for other cases.
)

// TODO: Add more fieldset for a map.
var UserFieldsetMap = map[UserFieldset][]string{
	UserFieldsetBasic: {
		"id",
		"name",
		"whatever",
	},
}
