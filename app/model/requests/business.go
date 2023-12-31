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
	Business_name *string            `json:"business_name"`
	Address 			*model.Address     `json:"address"`
	Latitude			*float64					 `json:"latitude"`
	Longitude			*float64					 `json:"longitude"`
	Description	  *string						 `json:"description"`	
}

type BusinessUserUpdate struct {
	First_name    *string  `json:"first_name"`
	Last_name     *string  `json:"last_name"`
	Password      *string  `json:"password"`
	Email         *string  `json:"email"`
	Business_name *string  `json:"business_name"`
	Address       *model.Address `json:"address"`
	Location      *model.Location `json:"location"`
	Description   *string   `json:"description"`	
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
		Business_name:   b.Business_name,
		Description:     b.Description,
		Password:				 b.Password,
		Email: 					 b.Email,
		Address:				b.Address,
	}
}