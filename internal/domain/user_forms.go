package domain

type UserWithForms struct {
	User  User   `json:"user"`
	Forms []Form `json:"forms"`
}
