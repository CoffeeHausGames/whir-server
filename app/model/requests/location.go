package model

import (
		"go.mongodb.org/mongo-driver/bson/primitive"
		"github.com/go-playground/validator/v10"
)


type Location struct {
	ID        primitive.ObjectID `bson:"_id"`
	Latitude  *float64           `json:"latitude" validate:"required,latitude"`
	Longitude *float64           `json:"longitude" validate:"required,longitude"`
	Radius    *float64           `json:"radius"`
}

// ValidateLocationStruct validates a Location struct
func ValidateLocationStruct(loc *Location) error {
	validate := validator.New()

	if err := validate.Struct(loc); err != nil {
		return err
	}

	return nil
}