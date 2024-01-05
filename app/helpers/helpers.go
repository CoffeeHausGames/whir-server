package helpers

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Contains(slice []*primitive.ObjectID, item *primitive.ObjectID) bool {
	for _, a := range slice {
			if a.Hex() == item.Hex() {
					return true
			}
	}
	return false
}