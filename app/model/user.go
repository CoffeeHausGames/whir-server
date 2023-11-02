package model

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

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

//ClientUser is the model that governs all account objects retrieved or inserted into the DB
type ClientUser struct {
	First_name    *string            `json:"first_name,omitempty"`
	Last_name     *string            `json:"last_name,omitempty"`
	Token         *string            `json:"token,omitempty"`
    Refresh_token *string            `json:"refresh_token,omitempty"`
    Base          *Base              `json:"Base,omitempty"`
}

// newUser sets up a client appropriate [model.User]
func NewUser(user *User) *ClientUser {
	return &ClientUser{
		First_name:      user.First_name,
		Last_name:       user.Last_name,
		Token:			 user.Token,
		Refresh_token:   user.Refresh_token,
	}
}