package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Address represents a structured address with various fields.
type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}

type Location struct {
	Type        	string    `json:"type" bson:"type"`
	Coordinates []float64 	`json:"coordinates" bson:"coordinates"`
}

//User is the model that governs all account objects retrieved or inserted into the DB
//TODO add back the bson
type BusinessUser struct {
		ID            primitive.ObjectID `bson:"_id"`
		First_name    *string            `json:"first_name" validate:"required,min=2,max=100"`
		Last_name     *string            `json:"last_name" validate:"required,min=2,max=100"`
		Password      *string            `json:"password" validate:"required,min=6"`
		Email         *string            `json:"email" validate:"email,required"`
		Token         *string            `json:"token"`
		Refresh_token *string            `json:"refresh_token"`
		Created_at    time.Time          `json:"created_at"`
		Updated_at    time.Time          `json:"updated_at"`
		Business_name *string            `json:"business_name"`
		Address 			*Address           `json:"address" bson:"address"`
		Location			*Location					 `json:"location" bson:"location"`
		Description	  *string						 `json:"description"`	
}

//BusinessUserWrapper is the model that represents the user to be sent to the frontend
type BusinessUserWrapper struct {
	First_name    *string            `json:"first_name,omitempty"`
	Last_name     *string            `json:"last_name,omitempty"`
	Business_name *string            `json:"business_name"`
	Address 			*Address           `json:"address"`
	Zip_code 			*string            `json:"zip_code"`
	Location			*Location					 `json:"location" bson:"location"`
	Deals 				[]*Deal	    			 `json:"deal"`	
	Description	  *string						 `json:"description"`	
}

// newUser sets up a frontend appropriate [model.User]
func NewBusinessAuthenticatedUser(business *BusinessUser, deals []*Deal) *BusinessUserWrapper {
	return &BusinessUserWrapper{
		First_name:      business.First_name,
		Last_name:       business.Last_name,
		Business_name:   business.Business_name,
		Address:			   business.Address,
		Location:				 business.Location,
		Deals: 					 deals,
		Description:     business.Description,
	}
}

// newUser sets up a frontend appropriate [model.User]
func NewBusinessUser(business *BusinessUser, deals []*Deal) *BusinessUserWrapper {
	return &BusinessUserWrapper{
		Business_name:   business.Business_name,
		Address:			   business.Address,
		Location:				 business.Location,
		Deals: 					 deals,
		Description:     business.Description,
	}
}

func GetStreetAddress(address *Address) string {
	// Concatenate the fields into a single street address string
	streetAddress := address.Street + ", " + address.City + ", " + address.State + " " + address.PostalCode + ", " + address.Country
	return streetAddress
}

func (b *BusinessUser) GetEmail() *string {
	return b.Email
}

func (b *BusinessUser) GetPassword() *string {
	return b.Password
}

func (b *BusinessUser) GetFirstName() *string {
	return b.First_name
}

func (b *BusinessUser) GetLastName() *string {
	return b.Last_name
}

func (b *BusinessUser) GetID() primitive.ObjectID {
	return b.ID
}

func (b *BusinessUser) GetToken() *string {
	return b.Token
}

func (b *BusinessUser) SetToken(token string) {
	b.Token = &token
}

func (b *BusinessUser) GetRefreshToken() *string {
	return b.Refresh_token
}

func (b *BusinessUser) SetRefreshToken(refreshToken string) {
	b.Refresh_token = &refreshToken
}