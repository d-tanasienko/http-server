package controller

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"httpserver/internal/config"
	"httpserver/internal/logger"
	"httpserver/internal/responses"
	"httpserver/internal/storage/activeuserstorage"
	"httpserver/internal/storage/tokenstorage"
	"httpserver/internal/storage/userstorage"
	"net/http"
	"time"
)

func UserHandler(writer http.ResponseWriter, request *http.Request, userStorage userstorage.UserStorageInterface, logger *logger.Logger) {
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

func UserLoginHandler(
	writer http.ResponseWriter,
	request *http.Request,
	userStorage userstorage.UserStorageInterface,
	logger *logger.Logger,
	tokenStorage tokenstorage.TokenStorageInterface,
) {
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

	token, err := generateSecureToken()
	if err != nil {
		logger.Error(err.Error())
		writer.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	tokenStorage.Add(token, user)

	url := "ws://" + config.GetBaseUrl() + config.GetPort() + "/ws?token=" + token
	responseData := responses.UserLoginResponse{Url: url}
	writer.Header().Add("X-Rate-Limit", "60")
	writer.Header().Add("X-Expires-After", currentTime.String())
	writer.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(writer)
	encoder.SetEscapeHTML(false)
	encoder.Encode(responseData)
}

func UserGetActiveList(w http.ResponseWriter, activeUsersStorage activeuserstorage.ActiveUsersStorageInterface) {
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.Encode(activeUsersStorage.GetNames())
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

func generateSecureToken() (string, error) {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
