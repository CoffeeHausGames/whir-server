package model

import (
	"time"
	"github.com/CoffeeHausGames/whir-server/app/model"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Deal struct {
	ID          primitive.ObjectID `json:"id,omitempty"`
	Business_id primitive.ObjectID `json:"business_id" bson:"business_id"`
	Name        *string            `json:"name"`
	Start_time  *time.Time         `json:"start_time"`
	End_time    *time.Time         `json:"end_time"`
	Day_of_week *string            `json:"day_of_week"`
	Start_date  *time.Time         `json:"start_date"`
	End_date  *time.Time         	 `json:"end_date"`
	Description *string            `json:"description"`
}

// ValidateLocationStruct validates a Location struct
func ValidateDealStruct(d *Deal) error {
	validate := validator.New()

	if err := validate.Struct(d); err != nil {
		return err
	}

	return nil
}

func NewDeal(d Deal) *model.Deal {
	return &model.Deal{
		ID: 					d.ID,
		Business_id:  d.Business_id,
		Name:      		d.Name,
		Start_time:		d.Start_time,
		End_time:   	d.End_time,
		Day_of_week:  d.Day_of_week,
		Start_date: 	d.Start_date,
		End_date:     d.End_date,
		Description:	d.Description,
	}
}