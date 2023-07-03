package controller

import (
	"context"
	"httpserver/internal/storage/activeuserstorage"
	"httpserver/internal/storage/tokenstorage"
	"net/http"

	"github.com/pkgz/websocket"
	"go.uber.org/zap"
)

func Ws(
	w http.ResponseWriter,
	r *http.Request,
	tokenStorage tokenstorage.TokenStorageInterface,
	activeUsersStorage activeuserstorage.ActiveUsersStorageInterface,
	logger *zap.SugaredLogger,
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
