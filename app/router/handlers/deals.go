package handlers

import (
		"encoding/json"
    "context"
		"log"
		"fmt"

    "net/http"
		"time"

		"github.com/julienschmidt/httprouter"

    "github.com/CoffeeHausGames/whir-server/app/auth"
		"github.com/CoffeeHausGames/whir-server/app/model"
		requests "github.com/CoffeeHausGames/whir-server/app/model/requests"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)


// Chat-GPTs attempt to make a Deal handler, help me :'(
	
func (env *HandlerEnv) AddDeal(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		claims := r.Context().Value("claims").(*auth.SignedDetails)
		body := r.Context().Value("body").(string)

    // Parse the request body to get building placement data
    // Unmarshal the URL-decoded JSON data into the placementData struct
    var dealData requests.Deal
    err := json.Unmarshal([]byte(body), &dealData)
    if err != nil {
        // Handle JSON decoding error
        WriteErrorResponse(w, http.StatusBadRequest, err.Error())
        return
    }
		
		err = requests.ValidateDealStruct(&dealData)
    if err != nil {
        // Handle JSON decoding error
        WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
        return
		}
		
		// Convert the hexadecimal string to an ObjectID
		fmt.Println(claims.Uid)
		objectID, err := primitive.ObjectIDFromHex(claims.Uid)
		if err != nil {
				// Handle the error if the hex string is not a valid ObjectID
				panic(err)
		}

		deal := requests.NewDeal(dealData)

		currBusiness := new(model.BusinessUser)
		var businessCollection model.Collection = env.database.GetBusinesses()
		err = businessCollection.FindOne(currBusiness, ctx, bson.M{"_id": objectID})

		currBusiness.Deals = append(currBusiness.Deals, deal)

    filter := bson.M{"_id": objectID} // Modify this filter according to your user ID format
		update := bson.M{"$set": bson.M{"deals": currBusiness.Deals}}

    _, err = businessCollection.UpdateOne(ctx, filter, update)
    if err != nil {
        log.Println(err)
        WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update user deals")
        return
    }

    WriteSuccessResponse(w, currBusiness)
}