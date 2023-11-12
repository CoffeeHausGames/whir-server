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

		"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)


// Function to add a deal for the authenticated business user
func (env *HandlerEnv) AddDeal(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		claims, body, err := env.getClaimsAndBody(r)
		if err != nil {
				WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
		}
	
		dealData, err := env.getDealDataFromBody(body)
		if err != nil {
				WriteErrorResponse(w, http.StatusBadRequest, err.Error())
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

	claims, body, err := env.getClaimsAndBody(r)
	if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
	}

	dealData, err := env.getDealDataFromBody(body)
	if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
	}

	userID, err := primitive.ObjectIDFromHex(claims.Uid)
	if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, "Error parsing user ID")
			return
	}

	dealCollection := env.database.GetDeals()
	deal := requests.NewDeal(dealData)
	deal.Business_id = userID

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

// function to delete deals for the authenticated business user
func (env *HandlerEnv) DeleteDeal(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	claims, body, err := env.getClaimsAndBody(r)
	if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
	}

	dealData, err := env.getDealDataFromBody(body)
	if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
	}

	userID, err := primitive.ObjectIDFromHex(claims.Uid)
	if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, "Error parsing user ID")
			return
	}

	dealCollection := env.database.GetDeals()
	deal := requests.NewDeal(dealData)
	deal.Business_id = userID

	query := bson.M{"_id": deal.ID, "business_id": userID}

	deleted, err := dealCollection.DeleteOne(ctx, query)
	if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, "Failed to delete deal")
			return
	}

	WriteSuccessResponse(w, deleted)
}


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

func (env *HandlerEnv) getClaimsAndBody(r *http.Request) (*auth.SignedDetails, string, error) {
	claims := r.Context().Value("claims").(*auth.SignedDetails)
	body := r.Context().Value("body").(string)
	return claims, body, nil
}

func (env *HandlerEnv) getDealDataFromBody(body string) (requests.Deal, error) {
	var dealData requests.Deal
	err := json.Unmarshal([]byte(body), &dealData)
	if err != nil {
			return dealData, err
	}

	err = requests.ValidateDealStruct(&dealData)
	if err != nil {
			return dealData, err
	}

	return dealData, nil
}