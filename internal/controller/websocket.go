package controller

import (
	"context"
	"httpserver/internal/logger"
	"httpserver/internal/storage/tokenstorage"
	"httpserver/internal/storage/userstorage"
	"net/http"

	"github.com/pkgz/websocket"
)

func Ws(
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