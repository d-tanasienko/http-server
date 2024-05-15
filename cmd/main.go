package main

import (
	"log"
	"net/http"

	"httpserver/internal/config"
	"httpserver/internal/controller"
	"httpserver/internal/storage/activeuserstorage"
	"httpserver/internal/storage/tokenstorage"
	"httpserver/internal/storage/userstorage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func main() {
	router := chi.NewRouter()
	userStorage := userstorage.NewUserStorage()
	tokenStorage := tokenstorage.NewTokenStorage()
	activeUsersStorage := activeuserstorage.NewActiveUsersStorage()
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Post("/user", func(w http.ResponseWriter, r *http.Request) {
		controller.UserHandler(w, r, userStorage, sugar)
	})

	router.Post("/user/login", func(w http.ResponseWriter, r *http.Request) {
		controller.UserLoginHandler(w, r, userStorage, sugar, tokenStorage)
	})

	router.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		controller.Ws(w, r, tokenStorage, activeUsersStorage, sugar)
	})
	router.Get("/user/active/list", func(w http.ResponseWriter, r *http.Request) {
		controller.UserGetActiveList(w, activeUsersStorage)
	})

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(config.GetPort(), nil))
}
