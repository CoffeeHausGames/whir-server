package handlers

import (
		"encoding/json"
    "context"
    "fmt"
		"log"

    "net/http"
		"time"

		"github.com/julienschmidt/httprouter"

    "github.com/CoffeeHausGames/whir-server/app/auth"
		"github.com/CoffeeHausGames/whir-server/app/model"
		"github.com/CoffeeHausGames/whir-server/app/helpers"
		requests "github.com/CoffeeHausGames/whir-server/app/model/requests"

		"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// Sign up allows a user with a unique email address to create an account and persists the account
// TODO see what code is repeated with login and make external function for it
func (env *HandlerEnv) BusinessSignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var businessCollection model.Collection = env.database.GetBusinesses()
	var userRequest requests.Business
	var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// pull the URL-decoded body from the context (comes from url_decoder middleware)
	decodedData := r.Context().Value("body").(string)

	err := json.Unmarshal([]byte(decodedData), &userRequest)
	if err != nil {
		log.Println(err)
		WriteErrorResponse(w, 422, "There was an error with the client request")
		return 
	}

	err = validate.Struct(userRequest)
	if err != nil {
		log.Println(err)
		WriteErrorResponse(w, 400, "There was an error with user validation")
		return
	}

	count, err := businessCollection.CountDocuments(ctx, bson.M{"email": userRequest.Email})
	defer cancel()
	if err != nil {
			log.Println(err)
			WriteErrorResponse(w, 502, "There was an error connecting with the server")
			return
	}

	if count > 0 {
		log.Println("error: this email already exists")
		// TODO update with email responses like described here
		// https://security.stackexchange.com/questions/51856/e-mail-already-in-use-exploit
		WriteErrorResponse(w, 401, "There was an error registering this account")
		return
	}
	user := requests.NewBusinessUser(userRequest)

	password := HashPassword(*user.Password)
	user.Password = &password

	//TODO sanitize input before inserting into DB to avoid NoSQL injection
	user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	token, refreshToken, _ := auth.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, user.ID.Hex())
	user.Token = &token
	user.Refresh_token = &refreshToken
	lon, lat := 0.0, 0.0
	if (userRequest.Longitude == nil || userRequest.Latitude == nil) && userRequest.Address != nil {
		address := model.GetStreetAddress(userRequest.Address)
		lon, lat, err = helpers.RetrieveCoordinatesFromAddress(address)
		fmt.Println(lon)
		fmt.Println(lat)
		if err != nil {
			log.Println("Address was not successful")
		}
	} else {
		lon = *userRequest.Longitude
		lat = *userRequest.Latitude
	}
	fmt.Println(lon)
	fmt.Println(lat)
	user.Location = &model.Location{
		Type: "Point",
		Coordinates: []float64{lon, lat},
	}

	_, insertErr := businessCollection.InsertOne(ctx, user)
	if insertErr != nil {
			log.Println("User item was not created")
			WriteErrorResponse(w, 502, "There was an error connecting with the server")
			return
	}
	defer cancel()

	WriteSuccessResponse(w, "Account created successfully")
}

// TODO move all the DB stuff to the model so we don't need to repeat code to get Users? Not sure what the golang standard here is 
//Login will allow a user to login to an account
func (env *HandlerEnv) BusinessLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var businessCollection model.Collection = env.database.GetBusinesses()
	var user model.BusinessUser
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// pull the URL-decoded body from the context (comes from url_decoder middleware)
	decodedData := r.Context().Value("body").(string)

	// Now you have the decoded JSON data as a string
	fmt.Println("Decoded JSON data:", decodedData)

	foundUser := new(model.BusinessUser)

	err := json.Unmarshal([]byte(decodedData), &user)
	if err != nil {
		log.Println(err)
		WriteErrorResponse(w, 422, "There was an error with the client request")
		return
	}

	err = validate.Var(user.Email, "required,email")
	if err != nil {
		log.Println(err)
		WriteErrorResponse(w, 400, "There was an error with user validation")
		return
	}

	//TODO sanitize input before Finding in DB to avoid NoSQL injection
	err = businessCollection.FindOne(foundUser, ctx, bson.M{"email": user.Email})
	defer cancel()
	if err != nil {
			log.Println(err)
			WriteErrorResponse(w, 502, "There was an error connecting with the server")
			return
	}

	passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
	defer cancel()
	if passwordIsValid != true {
		log.Println(msg)
		WriteErrorResponse(w, 401, "The username or password is incorrect")
		return
	}

	token, refreshToken, _ := auth.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, foundUser.ID.Hex())

	auth.UpdateAllTokens(businessCollection, token, refreshToken, foundUser.ID.Hex())

	userWrapper := model.NewBusinessUser(foundUser)
	userWrapper.Token = &token
	userWrapper.Refresh_token = &refreshToken

	WriteSuccessResponse(w, userWrapper)
}

func (env *HandlerEnv) BusinessTokenRefresh(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var businessCollection model.Collection = env.database.GetBusinesses()
	clientToken := r.Header.Get("refresh_token")
	if clientToken == "" {
			log.Printf("There is no refresh token")

			return
	}

	claims, err := auth.ValidateToken(businessCollection, clientToken)
	if err != "" {
			log.Panic(err)
			return
	}
	
	token, refreshToken, _ := auth.GenerateAllTokens(claims.Email, claims.First_name, claims.Last_name, claims.Uid)

	auth.UpdateAllTokens(businessCollection, token, refreshToken, claims.Uid)

	var user model.BusinessUserWrapper
	user.Refresh_token = &refreshToken
	user.Token = &token

	WriteSuccessResponse(w, user)
}

// Function to place a building in a base
//TODO make it so a user can send in an address and find the long and lat from that address
func (env *HandlerEnv) GetBusiness(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	radius := 20.0
	var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	body := r.Context().Value("body").(string)
	//TODO if lat and long is null then check address to find it

	// Parse the request body to get building placement data
	// Unmarshal the URL-decoded JSON data into the placementData struct
	var locationData requests.Location
	err := json.Unmarshal([]byte(body), &locationData)
	if err != nil {
			// Handle JSON decoding error
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
	}
	
	err = requests.ValidateLocationStruct(&locationData)
	if err != nil {
			// Handle JSON decoding error
			WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
			return
	}

	if locationData.Radius != nil {
		radius = *locationData.Radius
	}

	radiusMeters := radius * 1609.34
	var businessCollection model.Collection = env.database.GetBusinesses()

	// Create a query for geospatial location within the radius
	query := bson.M{
		"location": bson.M{
				"$near": bson.M{
						"$geometry": bson.M{
								"type":        "Point",
								"coordinates": []float64{*locationData.Longitude, *locationData.Latitude}, // The order is longitude (X), then latitude (Y).
						},
						"$maxDistance": radiusMeters, // Radius in meters.
				},
		},
	}

	// cursor, err := businessCollection.Find(currBusiness, ctx, bson.M{"zip_code": *locationData.Zip_code})
	cursor, err := businessCollection.Find(ctx, query);
	if err != nil {
		WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	businesses := GetMultipleBusinesses(cursor)

	// businessWrapper := model.NewBusinessUser(currBusiness)

	// Return a success response to the client
	WriteSuccessResponse(w, businesses)
}


// // Function to place a building in a base
// func (env *HandlerEnv) GetBusiness(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 	currBusiness := new(model.BusinessUser)
// 	var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	// How to access the "claims" object so the user properties from the auth token
// 	claims := r.Context().Value("claims").(*auth.SignedDetails)
// 	body := r.Context().Value("body").(string)

// 	// Parse the request body to get building placement data
// 	// Unmarshal the URL-decoded JSON data into the placementData struct
// 	var locationData requests.Location
// 	err := json.Unmarshal([]byte(body), &locationData)
// 	if err != nil {
// 			// Handle JSON decoding error
// 			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
// 			return
// 	}
	
// 	err = requests.ValidateLocationStruct(&locationData)
// 	if err != nil {
// 			// Handle JSON decoding error
// 			WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
// 			return
// 	}

// 	var businessCollection model.Collection = env.database.GetUsers()
// 	err = businessCollection.FindOne(currBusiness, ctx, bson.M{"zip_code": locationData.Zip_code})

// 	err = baseCollection.FindOne(currBase, ctx, bson.M{"owner": objectID})
// 	building := model.NewBuilding(&placementData)

// 	// Perform validation checks (e.g., user permissions, available resources)

// 	// Need to persist the base after the building is placed
// 	if err := currBase.AddBuildingToBase(building); err != nil {
// 			// Validation failed, return an error response
// 			WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
// 			return
// 	}

// 	// Define the filter to identify the base document by its ID (assuming it's a unique identifier).
// 	filter := bson.M{"_id": currBase.ID}

// 	// Define the update operation to set the "Grid" field in the base document.
// 	update := bson.M{
// 			"$set": bson.M{
// 					"Grid": currBase.Grid,
// 					"Buildings": currBase.Buildings,
// 					// Add any other fields you need to update here
// 			},
// 	}

// 	// Specify update options (optional). For example, you can enable upsert or specify additional options.
// 	options := options.Update().SetUpsert(false)

// 	// Perform the update operation on the "bases" collection.
// 	result, err := baseCollection.UpdateOne(ctx, filter, update, options)
// 	if err != nil {
// 			// Handle the error if the update operation fails
// 			WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update base: "+err.Error())
// 			return
// 	}

// 	if result.MatchedCount == 0 {
// 			// Handle the case where no documents were matched by the filter
// 			WriteErrorResponse(w, http.StatusNotFound, "No matching base found for the update")
// 			return
// 	}
	
// 	if result.ModifiedCount == 0 {
// 			// Handle the case where the update didn't modify any documents
// 			WriteErrorResponse(w, http.StatusNoContent, "No changes were made during the update")
// 			return
// 	}

// 	// Return a success response to the client
// 	WriteSuccessResponse(w, currBase)
// }



func GetMultipleBusinesses(cursor *mongo.Cursor) []model.BusinessUserWrapper{

	var businesses []model.BusinessUserWrapper
	// Iterate through the cursor and decode each document into a Businesses struct.
	for cursor.Next(context.Background()) {
		var business model.BusinessUserWrapper

		// Decode the current document into the Businesses struct.
		if err := cursor.Decode(&business); err != nil {
				log.Fatal(err)
		}

		// Append the decoded struct to the slice.
		businesses = append(businesses, business)
}
return businesses
}