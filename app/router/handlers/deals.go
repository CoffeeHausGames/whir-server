package handlers

import (
	"fmt"
		"encoding/json"
    "context"
		"log"

    "net/http"
		"time"

		"github.com/julienschmidt/httprouter"

    "github.com/CoffeeHausGames/whir-server/app/auth"
		"github.com/CoffeeHausGames/whir-server/app/model"
		requests "github.com/CoffeeHausGames/whir-server/app/model/requests"

		"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)


// Function to add a deal for the authenticated business user
//I need to edit the below method to use the new database methods and store the deals in its own collection with the business ID as a foreign key
	
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
		objectID, err := primitive.ObjectIDFromHex(claims.Uid)
		if err != nil {
				// Handle the error if the hex string is not a valid ObjectID
				panic(err)
		}
		fmt.Println(objectID)
		dealData.Business_id = objectID
		deal := requests.NewDeal(dealData)

    dealCollection := env.database.GetDeals()
    InsertedID, err := dealCollection.InsertOne(ctx, deal)
    if err != nil {
        log.Println(err)
        WriteErrorResponse(w, http.StatusInternalServerError, "Failed to insert deal")
        return
    }
		deal.ID = InsertedID

    WriteSuccessResponse(w, deal)
}
// TODO remove deals from businessUser model stored in bson and have it retrieve all the deals and return them from deal collection

// Function to get deals for the authenticated business user
func (env *HandlerEnv) GetSignedInBusinessDeals(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	claims := r.Context().Value("claims").(*auth.SignedDetails)
	userID, err := primitive.ObjectIDFromHex(claims.Uid)
	if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, "Error parsing user ID")
			return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	dealCollection := env.database.GetDeals()

	query := bson.M{"business_id": userID}

	cursor, err := dealCollection.Find(ctx, query)
	if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, "Error finding deals")
			return
	}

	deals := GetMultipleDeals(cursor)

	WriteSuccessResponse(w, deals)
}

// function to update deals for the authenticated user
func (env *HandlerEnv) UpdateDeal(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	claims := r.Context().Value("claims").(*auth.SignedDetails)
	body := r.Context().Value("body").(string)

	var dealData requests.Deal
	err := json.Unmarshal([]byte(body), &dealData)
	if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
	}

	err = requests.ValidateDealStruct(&dealData)
	if err != nil {
			WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
			return
	}

	userID, err := primitive.ObjectIDFromHex(claims.Uid)
	if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, "Error parsing user ID")
			return
	}

	dealCollection := env.database.GetDeals()
	deal := requests.NewDeal(dealData)

	query := bson.M{"_id": deal.ID, "business_id": userID}

	update := bson.M{
			"$set": deal,
	}

	_, err = dealCollection.UpdateOne(ctx, query, update)
	if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update deal")
			return
	}

	WriteSuccessResponse(w, deal)
}

// // Function to update a deal for the authenticated business user
// func (env *HandlerEnv) UpdateDeal(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 	var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	claims := r.Context().Value("claims").(*auth.SignedDetails)
// 	body := r.Context().Value("body").(string)

// 	// Need to retrieve the deal that is being updated from the body
// 	var dealData requests.Deal
// 	err := json.Unmarshal([]byte(body), &dealData)
// 	if err != nil {
// 			// Handle JSON decoding error
// 			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
// 			return
// 	}

// 	err = requests.ValidateDealStruct(&dealData)
// 	if err != nil {
// 			// Handle JSON decoding error
// 			WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
// 			return
// 	}

// 	// Convert the hexadecimal string to an ObjectID
// 	objectID, err := primitive.ObjectIDFromHex(claims.Uid)
// 	if err != nil {
// 			// Handle the error if the hex string is not a valid ObjectID
// 			WriteErrorResponse(w, http.StatusNotFound, "Business user not found")
// 	}

// 	// Assuming that the "business" collection contains the business users
// 	var businessCollection model.Collection = env.database.GetBusinesses()

// 	// Create a query to find the business user by ID
// 	query := bson.M{"_id": objectID}


// }

// Function to retrieve multiple businesses near a location
func GetMultipleDeals(cursor *mongo.Cursor) []model.Deal{

	var deals []model.Deal
	// Iterate through the cursor and decode each document into a Businesses struct.
	for cursor.Next(context.Background()) {
		var deal model.Deal

		// Decode the current document into the Businesses struct.
		if err := cursor.Decode(&deal); err != nil {
				log.Fatal(err)
		}

		// Append the decoded struct to the slice.
		deals = append(deals, deal)
}
return deals
}