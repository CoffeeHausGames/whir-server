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
	token, refreshToken, _ := auth.GenerateAllTokens(*user.Email, *user.First_name, user.Last_name, user.ID.Hex())
	user.Token = &token
	user.Refresh_token = &refreshToken

	_, insertErr := userCollection.InsertOne(ctx, user)
	if insertErr != nil {
			log.Println("User item was not created")
			WriteErrorResponse(w, 502, "There was an error connecting with the server")
			return
	}
	defer cancel()

	WriteSuccessResponse(w, r, "Account created successfully", nil, false)
}

//Login will allow a user to login to an account
func (env *HandlerEnv) UserLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params,) {
	var userCollection model.Collection = env.database.GetUsers()
	var user model.User
	foundUser := new(model.User)

	foundUserInterface, err := performLogin(r.Context(), userCollection, &user, foundUser)
	if err != nil {
		log.Println("performLogin failed")
		WriteErrorResponse(w, 401, "There was an error logging in")
		return
	}

	foundUser, ok := foundUserInterface.(*model.User)
	if !ok {
		log.Println("Wrong user type found")
		WriteErrorResponse(w, 401, "There was an error logging in")
		return
	}

	userWrapper := model.NewUser(foundUser)

	WriteSuccessResponse(w, r, userWrapper, foundUser, true)
}

func (env *HandlerEnv) TokenRefresh(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var userCollection model.Collection = env.database.GetUsers()
	clientToken := r.Header.Get("refresh_token")
	if clientToken == "" {
			log.Printf("There is no refresh token")

			return
	}

	claims, err := auth.ValidateToken(userCollection, clientToken)
	if err != nil {
			log.Panic(err)
			return
	}
	
	token, refreshToken, _ := auth.GenerateAllTokens(claims.Email, claims.First_name, &claims.Last_name, claims.Uid)

	auth.UpdateAllTokens(userCollection, token, refreshToken, claims.Uid)

	var user model.User
	user.Refresh_token = &refreshToken
	user.Token = &token

	WriteSuccessResponse(w, r, nil, &user, true)
}	

func (env *HandlerEnv) GetLoggedInUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	currUser := new(model.User)
	var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	claims := r.Context().Value("claims").(*auth.SignedDetails)

	var userCollection model.Collection = env.database.GetUsers()

	Id, err := primitive.ObjectIDFromHex(claims.Uid)
	err = userCollection.FindOne(currUser, ctx, bson.M{"_id": Id})

	if err != nil {
		WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	user := model.NewUser(currUser)

	// Return a success response to the client
	WriteSuccessResponse(w, r, user, currUser, true)
}

func (env *HandlerEnv) Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Clear the auth_token cookie
	http.SetCookie(w, &http.Cookie{
			Name:    "access_token",
			Value:   "",
			Path:    "/",
			Secure:  true,
			HttpOnly: true,
			MaxAge:  -1, // MaxAge<0 means delete cookie now
	})

	// Clear the refresh_token cookie
	http.SetCookie(w, &http.Cookie{
			Name:    "refresh_token",
			Value:   "",
			Path:    "/",
			Secure:  true,
			HttpOnly: true,
			MaxAge:  -1, // MaxAge<0 means delete cookie now
	})

   // Use WriteSuccessResponse to send the success message
	 WriteSuccessResponse(w, r, "Logged out successfully", nil, false)
}