package domain

type UserWithForm struct {
	User User `json:"user"`
	Form Form `json:"form"`
}
