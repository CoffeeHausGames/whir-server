package model

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Deal struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        *string            `json:"name"`
	Start_time  *time.Time         `json:"start_time" bson:"start_time,omitempty"`
	End_time    *time.Time         `json:"end_time" bson:"end_time,omitempty"`
	Day_of_week *string            `json:"day_of_week" bson:"day_of_week,omitempty"`
	Start_date  *time.Time         `json:"start_date" bson:"start_date,omitempty"`
	End_date  *time.Time         	 `json:"end_date" bson:"end_date,omitempty"`
	Description *string            `json:"description" bson:"description,omitempty"`
}