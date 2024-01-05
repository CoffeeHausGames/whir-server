package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"fmt"
	"errors"
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
		ID            primitive.ObjectID 		`bson:"_id"`
		First_name    *string            		`json:"first_name" validate:"required,min=2,max=100"`
		Last_name     *string            		`json:"last_name" validate:"required,min=2,max=100"`
		Password      *string            		`json:"password" validate:"required,min=6"`
		Email         *string            		`json:"email" validate:"email,required"`
		Token         *string            		`json:"token"`
		Refresh_token *string            		`json:"refresh_token"`
		Created_at    time.Time          		`json:"created_at"`
		Updated_at    time.Time          		`json:"updated_at"`
		Business_name *string            		`json:"business_name"`
		Address 			*Address           		`json:"address" bson:"address"`
		Location			*Location					 		`json:"location" bson:"location"`
		Description	  *string						 		`json:"description"`	
		PinnedDeals	 []*primitive.ObjectID	`json:"pinned_deals"`
}

//BusinessUserWrapper is the model that represents the user to be sent to the frontend
type BusinessUserWrapper struct {
	ID            primitive.ObjectID 		`json:"id"`
	First_name    *string            		`json:"first_name,omitempty"`
	Last_name     *string            		`json:"last_name,omitempty"`
	Business_name *string            		`json:"business_name"`
	Address 			*Address           		`json:"address"`
	Zip_code 			*string            		`json:"zip_code"`
	Location			*Location					 		`json:"location" bson:"location"`
	Deals 				[]*Deal	    			 		`json:"deals"`	
	Description	  *string						 		`json:"description"`	
	PinnedDeals	 []*primitive.ObjectID	`json:"pinned_deals"`
}

// newUser sets up a frontend appropriate [model.User]
func NewBusinessAuthenticatedUser(business *BusinessUser, deals []*Deal) *BusinessUserWrapper {
	return &BusinessUserWrapper{
		ID: 						 business.ID,
		First_name:      business.First_name,
		Last_name:       business.Last_name,
		Business_name:   business.Business_name,
		Address:			   business.Address,
		Location:				 business.Location,
		Deals: 					 deals,
		Description:     business.Description,
		PinnedDeals:		 business.PinnedDeals,
	}
}

// newUser sets up a frontend appropriate [model.User]
func NewBusinessUser(business *BusinessUser, deals []*Deal) *BusinessUserWrapper {
	return &BusinessUserWrapper{
		ID: 						 business.ID,
		Business_name:   business.Business_name,
		Address:			   business.Address,
		Location:				 business.Location,
		Deals: 					 deals,
		Description:     business.Description,
		PinnedDeals:		 business.PinnedDeals,
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

func (business *BusinessUser) PinDeal(dealID *primitive.ObjectID) error {
	maxPinnedDeals := 3
	if len(business.PinnedDeals) >= maxPinnedDeals {
		return fmt.Errorf("cannot pin more than %d deals", maxPinnedDeals)
	}
	business.PinnedDeals = append(business.PinnedDeals, dealID)
	return nil
}

func (business *BusinessUser) UnpinDeal(dealID *primitive.ObjectID) error {
	// Unpin the deal
	for i, deal := range business.PinnedDeals {
			if deal.Hex() == dealID.Hex() {
					// Remove the deal from the slice
					business.PinnedDeals = append(business.PinnedDeals[:i], business.PinnedDeals[i+1:]...)
					return nil
			}
	}
	return errors.New("deal not found in pinned deals")
}