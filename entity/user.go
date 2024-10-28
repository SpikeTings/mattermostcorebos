package entity

import "github.com/mattermost/mattermost/server/public/model"

type User struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password,omitempty"`
	Email     string `json:"email"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Position  string `json:"position"`
	Roles     string `json:"roles"`
	TeamNames string `json:"teamnames"`
}

func (u *User) GetMMUser() model.User {
	return model.User{Username: u.Username,
		Password:  u.Password,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Position:  u.Position,
		Roles:     u.Roles}
}

func (u User) GetUser(mmUser *model.User) User {
	return User{Username: mmUser.Username,
		Id:        mmUser.Id,
		Email:     mmUser.Email,
		FirstName: mmUser.FirstName,
		LastName:  mmUser.LastName,
		Position:  mmUser.Position,
		Roles:     mmUser.Roles}
}
