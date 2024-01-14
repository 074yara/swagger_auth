package main

import (
	"encoding/json"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"log"
	"net/http"
)

var (
	userIDCounter = 1
	users         = make(map[int]*User)
	jwtKey        = []byte(gofakeit.HackerPhrase())
	tokenAuth     = jwtauth.New("HS256", jwtKey, nil)
)

type User struct {
	ID       int    `json:"-"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// registerHandler register a new user
//
//	@Summary		new user handler
//	@Description	registers a new user with hashed password and adds it to storage in memory
//	@Tags			register
//	@Accept			json
//	@Produce		html
//	@Param			input	body	User	true	"User"
//	@Success		201		"User registered"
//	@Failure		400
//	@Failure		409	"User already exists"
//	@Router			/register [post]
func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if _, exists := findUserByName(user.Username); exists {
		http.Error(w, "User \""+user.Username+"\" already exists", http.StatusConflict)
		return
	}
	hashedPassword, err := HashPassword(user.Password)
	checkError(err)

	user.Password = hashedPassword
	user.ID = userIDCounter
	users[user.ID] = &user
	userIDCounter++

	log.Printf("user \"%v\" with ID %v registered", user.Username, user.ID)

	w.WriteHeader(http.StatusCreated)

}

// loginHandler user login
//
//	@Summary		new user handler
//	@Description	login user into system
//	@Tags			login
//	@Accept			json
//	@Produce		html
//	@Param			input	body		User	true	"User"
//	@Success		200		{string}	string	"token"
//	@Failure		400
//	@Failure		401	"Incorrect username or password"
//	@Failure		500
//	@Router			/login [post]
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	foundUser, exists := findUserByName(user.Username)
	if !exists {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if !CheckPasswordHash(user.Password, foundUser.Password) {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"username": foundUser.Username, "id": foundUser.ID})
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	log.Printf("user %v with ID %v succesfully logged in", foundUser.Username, foundUser.ID)

	_, err = w.Write([]byte("Bearer " + tokenString))
	if err != nil {
		http.Error(w, "Failed to write token", http.StatusInternalServerError)
		return
	}
}

func JwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := jwtauth.TokenFromHeader(r)
		if tokenString == "" {
			http.Error(w, "Not authorized", http.StatusForbidden)
			return
		}

		log.Println(tokenString)

		token, err := tokenAuth.Decode(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusForbidden)
			return
		}

		err = jwt.Validate(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func findUserByName(username string) (*User, bool) {
	for _, user := range users {
		if user.Username == username {
			return user, true
		}
	}
	return nil, false
}
