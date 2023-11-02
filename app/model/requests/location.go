package model

import (
		"go.mongodb.org/mongo-driver/bson/primitive"
		"github.com/go-playground/validator/v10"
)


type Location struct {
	ID          primitive.ObjectID `bson:"_id"`
	Zip_code    *string						 `json:"zip_code" validate:"required,min=2,max=100"`
} // TODO look into data structure probably best to use a date for these fields

// ValidateLocationStruct validates a Location struct
func ValidateLocationStruct(bp *Location) error {
	validate := validator.New()

	if err := validate.Struct(bp); err != nil {
		return err
	}

	return nil
}