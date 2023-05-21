package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"httpserver/logger"
	"httpserver/responses"
	"httpserver/userstorage"

	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()
	userStorage := userstorage.NewUserStorage()
	logger := logger.NewLogger()
	router.Post("/user", func(w http.ResponseWriter, r *http.Request) {
		userHandler(w, r, userStorage, logger)
	})

	router.Post("/user/login", func(w http.ResponseWriter, r *http.Request) {
		userLoginHandler(w, r, userStorage, logger)
	})
	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(GetPort(), nil))
}

func userHandler(writer http.ResponseWriter, request *http.Request, userStorage *userstorage.UserStorage, logger *logger.Logger) {
	userName, password, err := getUsernameAndPasswordFromBody(request)
	if err != nil {
		logger.Error(err.Error())
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}
	id := userStorage.Add(userName, password)
	responseData := &responses.UserResponse{Id: id, UserName: userName}

	writer.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(writer)
	encoder.Encode(responseData)
}

func userLoginHandler(writer http.ResponseWriter, request *http.Request, userStorage *userstorage.UserStorage, logger *logger.Logger) {
	userName, password, err := getUsernameAndPasswordFromBody(request)
	if err != nil {
		logger.Error(err.Error())
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := userStorage.Get(userName)
	if err != nil || user.Password != password {
		logger.Error(err.Error())
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	currentTime := time.Now().UTC()
	currentTime = currentTime.Add(time.Hour * 1)

	token, err := GenerateSecureToken()
	if err != nil {
		logger.Error(err.Error())
		writer.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	url := "ws://fancy-chat.io/ws&token=" + token
	responseData := responses.UserLoginResponse{Url: url}
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

	if len(userName) < 4 {
		return "", "", errors.New("username should be 4 chars or longer")
	}

	if len(password) < 8 {
		return "", "", errors.New("password should be 8 chars or longer")
	}

	if !userNameOk || !passwordOk {
		return "", "", errors.New("invalid username or password")
	}

	return userName, password, nil
}

func GenerateSecureToken() (string, error) {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	return ":" + port
}
