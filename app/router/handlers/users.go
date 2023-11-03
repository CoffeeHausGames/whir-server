package handlers

import (
		"encoding/json"
    "context"
    "fmt"
		"log"

    "net/http"
		"time"

		"github.com/go-playground/validator/v10"
		"github.com/julienschmidt/httprouter"

    "github.com/CoffeeHausGames/whir-server/app/auth"
    "github.com/CoffeeHausGames/whir-server/app/model"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

//HashPassword is used to encrypt the password before it is stored in the DB
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
			log.Panic(err)
	}

	return string(bytes)
}

//VerifyPassword checks the input password while verifying it with the passward in the DB.
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
			msg = fmt.Sprintf("The username or password is incorrect")
			check = false
	}

	return check, msg
}

// Sign up allows a user with a unique email address to create an account and persists the account
// TODO see what code is repeated with login and make external function for it
func (env *HandlerEnv) SignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var userCollection model.Collection = env.database.GetUsers()
	var user model.User
	var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// pull the URL-decoded body from the context (comes from url_decoder middleware)
	decodedData := r.Context().Value("body").(string)

	err := json.Unmarshal([]byte(decodedData), &user)
	if err != nil {
		log.Println(err)
		WriteErrorResponse(w, 422, "There was an error with the client request")
		return 
	}

	err = validate.Struct(user)
	if err != nil {
		log.Println(err)
		WriteErrorResponse(w, 400, "There was an error with user validation")
		return
	}

	count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
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

	password := HashPassword(*user.Password)
	user.Password = &password

	//TODO sanitize input before inserting into DB to avoid NoSQL injection
	user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	token, refreshToken, _ := auth.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, user.ID.Hex())
	user.Token = &token
	user.Refresh_token = &refreshToken

	_, insertErr := userCollection.InsertOne(ctx, user)
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
func (env *HandlerEnv) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var userCollection model.Collection = env.database.GetUsers()
	var user model.User
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// pull the URL-decoded body from the context (comes from url_decoder middleware)
	decodedData := r.Context().Value("body").(string)

	// Now you have the decoded JSON data as a string
	fmt.Println("Decoded JSON data:", decodedData)

	foundUser := new(model.User)

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
	err = userCollection.FindOne(foundUser, ctx, bson.M{"email": user.Email})
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

	auth.UpdateAllTokens(userCollection, token, refreshToken, foundUser.ID.Hex())

	userWrapper := model.NewUser(foundUser)
	userWrapper.Token = &token
	userWrapper.Refresh_token = &refreshToken

	WriteSuccessResponse(w, userWrapper)
}

func (env *HandlerEnv) TokenRefresh(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var userCollection model.Collection = env.database.GetUsers()
	clientToken := r.Header.Get("refresh_token")
	if clientToken == "" {
			log.Printf("There is no refresh token")

			return
	}

	claims, err := auth.ValidateToken(userCollection, clientToken)
	if err != "" {
			log.Panic(err)
			return
	}
	
	token, refreshToken, _ := auth.GenerateAllTokens(claims.Email, claims.First_name, claims.Last_name, claims.Uid)

	auth.UpdateAllTokens(userCollection, token, refreshToken, claims.Uid)

	var user model.UserWrapper
	user.Refresh_token = &refreshToken
	user.Token = &token

	WriteSuccessResponse(w, user)
}	

func (env *HandlerEnv) GetUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	currUser := new(model.User)
	var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	claims := r.Context().Value("claims").(*auth.SignedDetails)

	// Parse the request body to get building placement data
	// Unmarshal the URL-decoded JSON data into the placementData struct

	var userCollection model.Collection = env.database.GetUsers()
	// Id, err := primitive.ObjectIDFromHex(claims.Uid)

	err := userCollection.FindOne(currUser, ctx, bson.M{"ID": claims.Uid})

	// businessWrapper := model.NewBusinessUser(currBusiness)
	if err != nil {
		WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	user := model.NewUser(currUser)
	// Return a success response to the client
	WriteSuccessResponse(w, user)
}



