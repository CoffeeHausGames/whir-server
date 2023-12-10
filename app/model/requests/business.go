package model

import (
	"github.com/CoffeeHausGames/whir-server/app/model"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Business struct {
	ID            primitive.ObjectID `bson:"_id"`
	First_name    *string            `json:"first_name" validate:"required,min=2,max=100"`
	Last_name     *string            `json:"last_name" validate:"required,min=2,max=100"`
	Password      *string            `json:"password" validate:"required,min=6"`
	Email         *string            `json:"email" validate:"email,required"`
	Token         *string            `json:"token"`
	Refresh_token *string            `json:"refresh_token"`
	Business_name *string            `json:"business_name"`
	Address 			*model.Address     `json:"address"`
	Latitude			*float64					 `json:"latitude"`
	Longitude			*float64					 `json:"longitude"`
	Description	  *string						 `json:"description"`	
	Cookie_consent *bool             `json:"cookie_consent"`
}

// ValidateLocationStruct validates a Location struct
func ValidateBuildingStruct(b *Business) error {
	validate := validator.New()

	if err := validate.Struct(b); err != nil {
		return err
	}

	return nil
}

func NewBusinessUser(b Business) *model.BusinessUser {
	return &model.BusinessUser{
		ID: b.ID,
		First_name:      b.First_name,
		Last_name:       b.Last_name,
		Token:			 		 b.Token,
		Refresh_token:   b.Refresh_token,
		Business_name:   b.Business_name,
		Description:     b.Description,
		Password:				 b.Password,
		Email: 					 b.Email,
		Address:				b.Address,
		Cookie_consent:  b.Cookie_consent,
	}
}