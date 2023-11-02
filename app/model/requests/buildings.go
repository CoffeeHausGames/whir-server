package model

import (
	"github.com/go-playground/validator/v10"
)

// Building represents a simple building with location
type BuildingPlacement struct {
	Name   string  `json:"name,omitempty" bson:"name,omitempty" validate:"required,min=2,max=100"`
	PosX   *float64 `json:"posX" bson:"posX" validate:"required,numeric"`
	PosY   *float64 `json:"posY" bson:"posY" validate:"required,numeric"`
	PosZ   *float64 `json:"posZ" bson:"posZ" validate:"required,numeric"`
	Width  float64 `json:"width,omitempty" bson:"width,omitempty" validate:"required,gte=1"`
	Height float64 `json:"height,omitempty" bson:"height,omitempty" validate:"required,gte=1"`
}

// ValidateBuildingPlacement validates a BuildingPlacement struct
func ValidateBuildingPlacementStruct(bp *BuildingPlacement) error {
	validate := validator.New()

	if err := validate.Struct(bp); err != nil {
		return err
	}

	return nil
}
