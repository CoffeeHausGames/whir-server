package model

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

type UserInterface interface {
    GetEmail() *string
    GetPassword() *string
    GetFirstName() *string
    GetLastName() *string
    GetID() primitive.ObjectID
    GetToken() *string
    SetToken(token string)
    GetRefreshToken() *string
    SetRefreshToken(refreshToken string)
}

//User is the model that governs all account objects retrieved or inserted into the DB
//TODO add back the bson
type User struct {
    ID            primitive.ObjectID `bson:"_id"`
    First_name    *string            `json:"first_name" validate:"required,min=2,max=100"`
    Last_name     *string            `json:"last_name" validate:"required,min=2,max=100"`
    Password      *string            `json:"password" validate:"required,min=6"`
    Email         *string            `json:"email" validate:"email,required"`
    Token         *string            `json:"token"`
    Refresh_token *string            `json:"refresh_token"`
    Created_at    time.Time          `json:"created_at"`
    Updated_at    time.Time          `json:"updated_at"`
}

//UserWrapper is the model that represents the user to be sent to the frontend
type UserWrapper struct {
	First_name    *string            `json:"first_name,omitempty"`
	Last_name     *string            `json:"last_name,omitempty"`
    Email         *string            `json:"email"`
}


// newUser sets up a client appropriate [model.User]
func NewUser(user *User) *UserWrapper {
	return &UserWrapper{
		First_name:      user.First_name,
		Last_name:       user.Last_name,
        Email:           user.Email,
	}
}

func (u *User) GetEmail() *string {
    return u.Email
}

func (u *User) GetPassword() *string {
    return u.Password
}

func (u *User) GetFirstName() *string {
    return u.First_name
}

func (u *User) GetLastName() *string {
    return u.Last_name
}

func (u *User) GetID() primitive.ObjectID {
    return u.ID
}

func (u *User) GetToken() *string {
    return u.Token
}

func (u *User) SetToken(token string) {
    u.Token = &token
}

func (u *User) GetRefreshToken() *string {
    return u.Refresh_token
}

func (u *User) SetRefreshToken(refreshToken string) {
    u.Refresh_token = &refreshToken
}
