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
    GetCookieConsent() *bool
    SetCookieConsent(cookieConsent bool)
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
    Cookie_consent *bool             `json:"cookie_consent"`
}

//UserWrapper is the model that represents the user to be sent to the frontend
type UserWrapper struct {
	First_name    *string            `json:"first_name,omitempty"`
	Last_name     *string            `json:"last_name,omitempty"`
    // TODO remove these so they are only sent in the header or in cookies
	Token         *string            `json:"token,omitempty"`
    Refresh_token *string            `json:"refresh_token,omitempty"`
}


// newUser sets up a client appropriate [model.User]
func NewUser(user *User) *UserWrapper {
	return &UserWrapper{
		First_name:      user.First_name,
		Last_name:       user.Last_name,
		Token:			 user.Token,
		Refresh_token:   user.Refresh_token,
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

// GetCookieConsent returns the user's cookie consent.
func (u *User) GetCookieConsent() *bool {
    return u.Cookie_consent
}

// SetCookieConsent sets the user's cookie consent.
func (u *User) SetCookieConsent(cookieConsent bool) {
    u.Cookie_consent = &cookieConsent
}