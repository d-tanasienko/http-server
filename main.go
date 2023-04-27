package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"httpserver/userstorage"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	userStorage := userstorage.NewUserStorage()
	router.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		userHandler(w, r, userStorage)
	}).Methods("POST")

	router.HandleFunc("/user/login", func(w http.ResponseWriter, r *http.Request) {
		userLoginHandler(w, r, userStorage)
	}).Methods("POST")
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func userHandler(writer http.ResponseWriter, request *http.Request, userStorage *userstorage.UserStorage) {
	userName, password, err := getUsernameAndPasswordFromBody(request)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
	} else {
		id := userStorage.Add(userName, password)
		responseData := map[string]string{
			"id":       id,
			"userName": userName,
		}
		writer.WriteHeader(http.StatusCreated)
		encoder := json.NewEncoder(writer)
		encoder.Encode(responseData)
	}
}

func userLoginHandler(writer http.ResponseWriter, request *http.Request, userStorage *userstorage.UserStorage) {
	userName, password, err := getUsernameAndPasswordFromBody(request)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := userStorage.Get(userName)
	if err != nil || user.Password != password {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	currentTime := time.Now().UTC()
	currentTime = currentTime.Add(time.Hour * 1)

	token := GenerateSecureToken()
	if token == "" {
		writer.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	url := "ws://fancy-chat.io/ws&token=" + token
	responseData := map[string]string{
		"url": url,
	}
	writer.Header().Add("X-Rate-Limit", "60")
	writer.Header().Add("X-Expires-After", currentTime.String())
	writer.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(writer)
	encoder.SetEscapeHTML(false)
	encoder.Encode(responseData)
}

func getUsernameAndPasswordFromBody(request *http.Request) (string, string, error) {
	decoder := json.NewDecoder(request.Body)
	var body = make(map[string]string)
	err := decoder.Decode(&body)
	if err != nil {
		return "", "", errors.New("invalid body")
	}

	userName, userNameOk := body["userName"]
	password, passwordOk := body["password"]

	if !userNameOk || !passwordOk || len(userName) < 4 || len(password) < 8 {
		return "", "", errors.New("invalid username or password")
	}

	return userName, password, nil
}

func GenerateSecureToken() string {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
