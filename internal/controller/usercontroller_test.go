package controller_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"

	"httpserver/internal/controller"
	"httpserver/internal/storage"
	"httpserver/internal/storage/tokenstorage"
)

type UserStorageMock struct {
}

func (m UserStorageMock) Add(userName string, password string) string {
	return "mocked_id"
}

func (m UserStorageMock) Get(userName string) (*storage.User, error) {
	return &storage.User{UserName: "JohnDoe", Password: "password123", Uuid: "mocked_id"}, nil
}

func TestUserHandler(t *testing.T) {
	userStorage := new(UserStorageMock)

	logger := zaptest.NewLogger(t).Sugar()

	reqBody := `{"userName": "JohnDoe", "password": "password123"}`
	req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	controller.UserHandler(w, req, userStorage, logger)

	assert.Equal(t, http.StatusCreated, w.Code)

	expectedResBody := `{"id":"mocked_id","userName":"JohnDoe"}`
	assert.Equal(t, expectedResBody, strings.TrimSpace(w.Body.String()))

	user, err := userStorage.Get("JohnDoe")
	assert.NoError(t, err)
	assert.Equal(t, "JohnDoe", user.UserName)
	assert.Equal(t, "password123", user.Password)
}

func TestUserLoginHandler(t *testing.T) {
	userStorage := new(UserStorageMock)
	logger := zaptest.NewLogger(t).Sugar()
	tokenStorage := tokenstorage.NewTokenStorage()

	reqBody := `{"userName": "JohnDoe","password": "password123"}`
	req := httptest.NewRequest(http.MethodPost, "/user/login", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	controller.UserLoginHandler(w, req, userStorage, logger, tokenStorage)

	assert.Equal(t, http.StatusCreated, w.Code)

	expectedURL := "ws://localhost:3000/ws?token="
	assert.Contains(t, w.Body.String(), expectedURL)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	token := strings.TrimSpace(strings.TrimPrefix(response[`url`], expectedURL))

	user, err := tokenStorage.Get(token)
	assert.NoError(t, err)
	assert.Equal(t, "JohnDoe", user.UserName)
	assert.Equal(t, "password123", user.Password)
}

func TestUserLoginHandler_InvalidBody(t *testing.T) {
	userStorage := new(UserStorageMock)
	buf := &zaptest.Buffer{}

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(buf),
		zapcore.ErrorLevel,
	))

	tokenStorage := tokenstorage.NewTokenStorage()

	reqBody := `{"userName": "JohnDoe","password": "password12`
	req := httptest.NewRequest(http.MethodPost, "/user/login", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	controller.UserLoginHandler(w, req, userStorage, logger.Sugar(), tokenStorage)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	logs := buf.String()
	assert.Contains(t, logs, "invalid body")
}

func TestUserLoginHandler_ShortUserName(t *testing.T) {
	userStorage := new(UserStorageMock)
	buf := &zaptest.Buffer{}

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(buf),
		zapcore.ErrorLevel,
	))

	tokenStorage := tokenstorage.NewTokenStorage()

	reqBody := `{"userName": "Jon","password": "password123"}`
	req := httptest.NewRequest(http.MethodPost, "/user/login", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	controller.UserLoginHandler(w, req, userStorage, logger.Sugar(), tokenStorage)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	logs := buf.String()
	assert.Contains(t, logs, "username should be 4 chars or longer")
}

func TestUserLoginHandler_ShortPassword(t *testing.T) {
	userStorage := new(UserStorageMock)
	buf := &zaptest.Buffer{}

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(buf),
		zapcore.ErrorLevel,
	))

	tokenStorage := tokenstorage.NewTokenStorage()

	reqBody := `{"userName": "John","password": "pass"}`
	req := httptest.NewRequest(http.MethodPost, "/user/login", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	controller.UserLoginHandler(w, req, userStorage, logger.Sugar(), tokenStorage)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	logs := buf.String()
	assert.Contains(t, logs, "password should be 8 chars or longer")
}

func TestUserLoginHandler_UserNameNotProvided(t *testing.T) {
	userStorage := new(UserStorageMock)
	buf := &zaptest.Buffer{}

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(buf),
		zapcore.ErrorLevel,
	))

	tokenStorage := tokenstorage.NewTokenStorage()

	reqBody := `{"password": "pass"}`
	req := httptest.NewRequest(http.MethodPost, "/user/login", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	controller.UserLoginHandler(w, req, userStorage, logger.Sugar(), tokenStorage)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	logs := buf.String()
	assert.Contains(t, logs, "invalid username or password")
}

func TestUserLoginHandler_PasswordNotProvided(t *testing.T) {
	userStorage := new(UserStorageMock)
	buf := &zaptest.Buffer{}

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(buf),
		zapcore.ErrorLevel,
	))

	tokenStorage := tokenstorage.NewTokenStorage()

	reqBody := `{"userName": "JohnDoe"}`
	req := httptest.NewRequest(http.MethodPost, "/user/login", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	controller.UserLoginHandler(w, req, userStorage, logger.Sugar(), tokenStorage)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	logs := buf.String()
	assert.Contains(t, logs, "invalid username or password")
}

type UserStorageInvalidUserMock struct {
}

func (m UserStorageInvalidUserMock) Add(userName string, password string) string {
	return "mocked_id"
}

func (m UserStorageInvalidUserMock) Get(userName string) (*storage.User, error) {
	return nil, errors.New("mocked error")
}

func TestUserLoginHandler_UserDoesNotExist(t *testing.T) {
	userStorage := new(UserStorageInvalidUserMock)
	buf := &zaptest.Buffer{}

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(buf),
		zapcore.ErrorLevel,
	))

	tokenStorage := tokenstorage.NewTokenStorage()

	reqBody := `{"userName": "John","password": "password123"}`
	req := httptest.NewRequest(http.MethodPost, "/user/login", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	controller.UserLoginHandler(w, req, userStorage, logger.Sugar(), tokenStorage)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	logs := buf.String()
	assert.Contains(t, logs, "mocked error")
}

type fakeActiveUsersStorage struct{}

func (m *fakeActiveUsersStorage) GetNames() []string {
	return []string{"User1", "User2", "User3"}
}

func (m *fakeActiveUsersStorage) Add(user *storage.User) {
}

func (m *fakeActiveUsersStorage) Get(userName string) (*storage.User, error) {
	return nil, nil
}

func (m *fakeActiveUsersStorage) Delete(user *storage.User) {
}

func TestUserGetActiveList(t *testing.T) {
	activeUsersStorage := &fakeActiveUsersStorage{}

	_, err := http.NewRequest("GET", "/users/active", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	controller.UserGetActiveList(rr, activeUsersStorage)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %d but got %d", http.StatusOK, rr.Code)
	}

	var names []string
	err = json.Unmarshal(rr.Body.Bytes(), &names)
	if err != nil {
		t.Errorf("failed to unmarshal response body: %v", err)
	}

	expectedLen := 3
	if len(names) != expectedLen {
		t.Errorf("expected %d active users but got %d", expectedLen, len(names))
	}
}
