package main

import (
	"encoding/json"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-chi/jwtauth/v5"
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
	Token    string `json:"-"`
}

// registerHandler register a new user
//
// @Summary new user handler
// @Description registers a new user with hashed password and adds it to storage in memory
// @Tags register
// @Accept json
// @Produce html
// @Param input body User true "User"
// @Router /register [post]
func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if _, exists := findUserByName(user.Username); exists {
		http.Error(w, "User `"+user.Username+"` already exists", http.StatusBadRequest)
		return
	}
	hashedPassword, err := HashPassword(user.Password)
	checkError(err)

	user.Password = hashedPassword
	user.ID = userIDCounter
	users[user.ID] = &user
	userIDCounter++

	log.Printf("user %v with ID %v registered", user.Username, user.ID)

	w.WriteHeader(http.StatusCreated)

}

// loginHandler user login
//
// @Summary new user handler
// @Description login user into system
// @Tags login
// @Accept json
// @Produce html
// @Param input body User true "User"
// @Router /login [post]
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

	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"user_id": foundUser.ID, "username": foundUser.Username})
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	users[user.ID].Token = tokenString

	log.Printf("user %v with ID %v succesfully logged in", foundUser.Username, foundUser.ID)

	_, err = w.Write([]byte(tokenString))
	if err != nil {
		http.Error(w, "Failed to write token", http.StatusInternalServerError)
		return
	}
}

func JwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			http.Error(w, "Authorization error", http.StatusUnauthorized)
			return
		}

		log.Println(claims)

		userID, ok := claims["user_id"].(int)
		if !ok {
			http.Error(w, "Unauthorized1", http.StatusUnauthorized)
			return
		}
		if _, exists := users[userID]; !exists {
			http.Error(w, "Unauthorized2", http.StatusUnauthorized)
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
