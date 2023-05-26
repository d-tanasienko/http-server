package main

import (
	"context"
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
	"httpserver/storage/tokenstorage"
	"httpserver/storage/userstorage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkgz/websocket"
)

func main() {
	router := chi.NewRouter()
	userStorage := userstorage.NewUserStorage()
	tokenStorage := tokenstorage.NewTokenStorage()
	activeUsersStorage := userstorage.NewActiveUsersStorage()
	logger := logger.NewLogger()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Post("/user", func(w http.ResponseWriter, r *http.Request) {
		userHandler(w, r, userStorage, logger)
	})

	router.Post("/user/login", func(w http.ResponseWriter, r *http.Request) {
		userLoginHandler(w, r, userStorage, logger, tokenStorage)
	})

	router.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws(w, r, tokenStorage, activeUsersStorage, logger)
	})
	router.Get("/user/active/list", func(w http.ResponseWriter, r *http.Request) {
		userGetActiveList(w, activeUsersStorage)
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

func userLoginHandler(
	writer http.ResponseWriter,
	request *http.Request,
	userStorage *userstorage.UserStorage,
	logger *logger.Logger,
	tokenStorage *tokenstorage.TokenStorage,
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

	token, err := GenerateSecureToken()
	if err != nil {
		logger.Error(err.Error())
		writer.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	tokenStorage.Add(token, user)

	url := "ws://localhost" + GetPort() + "/ws?token=" + token
	responseData := responses.UserLoginResponse{Url: url}
	writer.Header().Add("X-Rate-Limit", "60")
	writer.Header().Add("X-Expires-After", currentTime.String())
	writer.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(writer)
	encoder.SetEscapeHTML(false)
	encoder.Encode(responseData)
}

func userGetActiveList(w http.ResponseWriter, activeUsersStorage *userstorage.ActiveUsersStorage) {
	userNames := make([]string, len(*activeUsersStorage))

	i := 0
	for k := range *activeUsersStorage {
		userNames[i] = k
		i++
	}
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.Encode(userNames)
}

func ws(
	w http.ResponseWriter,
	r *http.Request,
	tokenStorage *tokenstorage.TokenStorage,
	activeUsersStorage *userstorage.ActiveUsersStorage,
	logger *logger.Logger,
) {
	user, err := tokenStorage.Get(r.URL.Query().Get("token"))
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	wsServer := websocket.Start(context.Background())
	activeUsersStorage.Add(user)
	wsServer.Handler(w, r)
	wsServer.On("echo", func(c *websocket.Conn, msg *websocket.Message) {
		err = c.Emit("echo", msg.Data)
		if err != nil {
			logger.Error(err.Error())
		}
	})
	activeUsersStorage.Delete(user)
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
