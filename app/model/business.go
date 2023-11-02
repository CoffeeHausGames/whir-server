package model

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)


type Deal struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        *string						 `json:"name"`
	Start_time  *string						 `json:"start_time"`
	End_time    *string            `json:"end_time"`
	Day_of_week	*string						 `json:"day_of_week"`
	Run_time		*string						 `json:"run_time"`
	Description	*string						 `json:"description"`
} // TODO look into data structure probably best to use a date for these fields


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
		Address 			*string            `json:"address"`
		Zip_code 			*string            `json:"zip_code"`
		Deals 				[]*Deal	    			 `json:"deal"`	
		Description	  *string						 `json:"description"`	
}

//BusinessUserWrapper is the model that represents the user to be sent to the frontend
type BusinessUserWrapper struct {
	First_name    *string            `json:"first_name,omitempty"`
	Last_name     *string            `json:"last_name,omitempty"`
	Token         *string            `json:"token,omitempty"`
	Refresh_token *string            `json:"refresh_token,omitempty"`
	Business_name *string            `json:"business_name"`
	Address 			*string            `json:"address"`
	Zip_code 			*string            `json:"zip_code"`
	Deals 				[]*Deal	    			 `json:"deal"`	
	Description	  *string						 `json:"description"`	
}

// newUser sets up a frontend appropriate [model.User]
func NewBusinessUser(user *BusinessUser) *BusinessUserWrapper {
	return &BusinessUserWrapper{
		First_name:      user.First_name,
		Last_name:       user.Last_name,
		Token:			 		user.Token,
		Refresh_token:   user.Refresh_token,
		Business_name: user.Business_name,
    Zip_code: user.Zip_code,
		Deals: user.Deals,
		Description: user.Description,
	}
}