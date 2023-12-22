package handlers

import (
	"reflect"
		"encoding/json"
    "context"
		"log"
		"errors"

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

	password := HashPassword(*user.GetPassword())
	user.Password = &password

	//TODO sanitize input before inserting into DB to avoid NoSQL injection
	user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	token, refreshToken, _ := auth.GenerateAllTokens(*user.GetEmail(), *user.GetFirstName(), user.GetLastName(), user.GetID().Hex())
	user.SetToken(token)
	user.SetRefreshToken(refreshToken)
	lon, lat := 0.0, 0.0
	if userRequest.Address != nil {
		address := model.GetStreetAddress(userRequest.Address)
		lon, lat, err = helpers.RetrieveCoordinatesFromAddress(address)
		if err != nil {
			log.Println("Address was not successful")
		}
	} else if (userRequest.Longitude != nil && userRequest.Latitude != nil) {
		lon = *userRequest.Longitude
		lat = *userRequest.Latitude
	}
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

	WriteSuccessResponse(w, r, "Account created successfully", nil, false)
}

// TODO move all the DB stuff to the model so we don't need to repeat code to get Users? Not sure what the golang standard here is 
//Login will allow a user to login to an account
func (env *HandlerEnv) BusinessLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var businessCollection model.Collection = env.database.GetBusinesses()
	var user model.BusinessUser
	foundUser := new(model.BusinessUser)

	foundUserInterface, err := performLogin(r.Context(), businessCollection, &user, foundUser)
	if err != nil {
		log.Println("performLogin failed")
		WriteErrorResponse(w, 401, "There was an error logging in")
		return
	}

	foundUser, ok := foundUserInterface.(*model.BusinessUser)
	if !ok {
		log.Println("Wrong user type found")
		WriteErrorResponse(w, 401, "There was an error logging in")
		return
	}

	userWrapper := model.NewBusinessAuthenticatedUser(foundUser, nil)

	WriteSuccessResponse(w, r, userWrapper, foundUser, true)
}

func (env *HandlerEnv) BusinessTokenRefresh(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var businessCollection model.Collection = env.database.GetBusinesses()
	clientToken := r.Header.Get("refresh_token")
	if clientToken == "" {
			log.Printf("There is no refresh token")

			return
	}

	claims, err := auth.ValidateToken(businessCollection, clientToken)
	if err != nil {
			log.Panic(err)
			return
	}
	
	token, refreshToken, _ := auth.GenerateAllTokens(claims.Email, claims.First_name, &claims.Last_name, claims.Uid)

	auth.UpdateAllTokens(businessCollection, token, refreshToken, claims.Uid)

	var user model.BusinessUser
	user.Refresh_token = &refreshToken
	user.Token = &token

	WriteSuccessResponse(w, r, nil, &user, false)
}

// Function to retrieve all businesses with deals
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
		WriteErrorResponse(w, http.StatusUnprocessableEntity, "Error retrieving businesses")
		return
	}

	// Decode the businesses
	var businesses []*model.BusinessUser
	if err := cursor.All(ctx, &businesses); err != nil {
		WriteErrorResponse(w, http.StatusUnprocessableEntity, "Error retrieving businesses")
		return
	}

	dealCollection := env.database.GetDeals()
	businessUserWrappers := make([]model.BusinessUserWrapper, 0, len(businesses))
	for _, business := range businesses {
		deals, err := GetDealForBusiness(business.ID, dealCollection, context.Background())
		if err != nil {
			WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
			continue
		}
		businessUserWrapper := model.NewBusinessUser(business, deals)

		businessUserWrappers = append(businessUserWrappers, *businessUserWrapper)
	}

	// Return a success response to the client
	WriteSuccessResponse(w, r, businessUserWrappers, nil, false)
}

// Function to retrieve multiple businesses near a location
func GetMultipleBusinesses(cursor *mongo.Cursor) ([]model.BusinessUser, error){

	var businesses []model.BusinessUser
	// Iterate through the cursor and decode each document into a Businesses struct.
	for cursor.Next(context.Background()) {
		var business model.BusinessUser

		// Decode the current document into the Businesses struct.
		if err := cursor.Decode(&business); err != nil {
				return nil, errors.New("Error decoding business")
		}

		// Append the decoded struct to the slice.
		businesses = append(businesses, business)
	}
	return businesses, nil
}

// Used to get a business by a given ID does not require auth and will return a business with deals to the client
func (env *HandlerEnv) GetBusinessByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	business := new(model.BusinessUser)
	var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var businessCollection model.Collection = env.database.GetBusinesses()

	Id, err := primitive.ObjectIDFromHex(ps.ByName("id"))
	err = businessCollection.FindOne(business, ctx, bson.M{"_id": Id})

	if err != nil {
		WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	dealCollection := env.database.GetDeals()

	deals, err := GetDealForBusiness(business.ID, dealCollection, context.Background())
	if err != nil {
		WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	businessUserWrapper := model.NewBusinessUser(business, deals)

	// Return a success response to the client
	WriteSuccessResponse(w, r, businessUserWrapper, nil, false)
}

// Updates the Business User Information
// NOTE this will not update deals. Deals are updated in the deals.go file
func (env *HandlerEnv) UpdateBusinessInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var businessCollection model.Collection = env.database.GetBusinesses()
	var userRequest requests.Business
	var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	currBusiness:= new(requests.Business)

	claims := r.Context().Value("claims").(*auth.SignedDetails)

	Id, err := primitive.ObjectIDFromHex(claims.Uid)
	err = businessCollection.FindOne(currBusiness, ctx, bson.M{"_id": Id})

	// pull the URL-decoded body from the context (comes from url_decoder middleware)
	decodedData := r.Context().Value("body").(string)

	err = json.Unmarshal([]byte(decodedData), &userRequest)
	if err != nil {
			log.Println(err)
			WriteErrorResponse(w, 422, "There was an error with the client request")
			return 
	}
	businessUser := requests.NewBusinessUser(userRequest)
	
	update := bson.M{}
	val := reflect.ValueOf(businessUser)
	typ := val.Type()
	
	for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			if !field.IsNil() {
					update[typ.Field(i).Name] = field.Interface()
			}
	}
	
	updateOperation := bson.M{
			"$set": update,
	}

	_, err = businessCollection.UpdateOne(ctx, bson.M{"_id": Id}, updateOperation)
	if err != nil {
		log.Println(err)
		WriteErrorResponse(w, 502, "There was an error updating the business info")
		return
	}

	WriteSuccessResponse(w, r, "Business info updated successfully", nil, false)
}

func (env *HandlerEnv) GetLoggedInBusiness(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	currBusiness:= new(model.BusinessUser)
	var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	claims := r.Context().Value("claims").(*auth.SignedDetails)

	var businessCollection model.Collection = env.database.GetBusinesses()

	Id, err := primitive.ObjectIDFromHex(claims.Uid)
	err = businessCollection.FindOne(currBusiness, ctx, bson.M{"_id": Id})

	if err != nil {
		WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	dealCollection := env.database.GetDeals()

	deals, err := GetDealForBusiness(currBusiness.ID, dealCollection, context.Background())
	if err != nil {
		WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	businessUserWrapper := model.NewBusinessUser(currBusiness, deals)

	// Return a success response to the client
	WriteSuccessResponse(w, r, businessUserWrapper, currBusiness, true)
}
