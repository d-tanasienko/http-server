package controller_test

import (
	"errors"
	"httpserver/internal/controller"
	"httpserver/internal/storage"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

type mockTokenStorage struct{}

func (m *mockTokenStorage) Get(token string) (*storage.User, error) {
	if token == "valid_token" {
		return &storage.User{UserName: "JohnDoe"}, nil
	}
	return nil, errors.New("invalid token")
}
func (m *mockTokenStorage) Add(string, *storage.User) {

}

func (m *mockTokenStorage) Delete(token string) {

}

type mockActiveUsersStorage struct {
	addedUser   *storage.User
	deletedUser *storage.User
}

func (m *mockActiveUsersStorage) Add(user *storage.User) {
	m.addedUser = user
}

func (m *mockActiveUsersStorage) Get(userName string) (*storage.User, error) {
	return nil, nil
}

func (m *mockActiveUsersStorage) GetNames() []string {
	return nil
}

func (m *mockActiveUsersStorage) Delete(user *storage.User) {
	m.deletedUser = user
}

func TestWs_ValidToken(t *testing.T) {
	tokenStorage := &mockTokenStorage{}
	activeUsersStorage := &mockActiveUsersStorage{
		addedUser:   &storage.User{},
		deletedUser: &storage.User{},
	}
	logger := zaptest.NewLogger(t).Sugar()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controller.Ws(w, r, tokenStorage, activeUsersStorage, logger)
	}))
	defer server.Close()

	dialer := websocket.Dialer{}

	url := "ws" + server.URL[4:] + "/ws?token=valid_token"

	conn, _, err := dialer.Dial(url, nil)
	assert.NoError(t, err)

	conn.Close()

	assert.NotNil(t, activeUsersStorage.addedUser)
	assert.Equal(t, "JohnDoe", activeUsersStorage.addedUser.UserName)

	time.Sleep(100 * time.Millisecond)

	assert.NotNil(t, activeUsersStorage.deletedUser)
	assert.Equal(t, "JohnDoe", activeUsersStorage.deletedUser.UserName)
}

func TestWs_InvalidToken(t *testing.T) {
	buf := &zaptest.Buffer{}

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(buf),
		zapcore.ErrorLevel,
	))

	fakeTokenStorage := &mockTokenStorage{}
	fakeActiveUsersStorage := &mockActiveUsersStorage{
		addedUser:   &storage.User{},
		deletedUser: &storage.User{},
	}

	req := httptest.NewRequest("GET", "/ws?token=invalid_token", nil)
	w := httptest.NewRecorder()

	controller.Ws(w, req, fakeTokenStorage, fakeActiveUsersStorage, logger.Sugar())

	logs := buf.String()
	assert.Contains(t, logs, "invalid token")

	assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
}
