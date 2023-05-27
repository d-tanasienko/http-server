package main

import (
	"log"
	"net/http"

	"httpserver/internal/config"
	"httpserver/internal/controller"
	"httpserver/internal/logger"
	"httpserver/internal/storage/tokenstorage"
	"httpserver/internal/storage/userstorage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
		controller.UserHandler(w, r, userStorage, logger)
	})

	router.Post("/user/login", func(w http.ResponseWriter, r *http.Request) {
		controller.UserLoginHandler(w, r, userStorage, logger, tokenStorage)
	})

	router.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		controller.Ws(w, r, tokenStorage, activeUsersStorage, logger)
	})
	router.Get("/user/active/list", func(w http.ResponseWriter, r *http.Request) {
		controller.UserGetActiveList(w, activeUsersStorage)
	})

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(config.GetPort(), nil))
}
