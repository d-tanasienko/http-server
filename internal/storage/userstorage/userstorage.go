package userstorage

import (
	"errors"
	"httpserver/internal/storage"

	"github.com/google/uuid"
)

type UserStorageInterface interface {
	Add(string, string) string
	Get(string) (*storage.User, error)
}

type UserStorage map[string]*storage.User

func (userStorage UserStorage) Add(userName string, password string) string {
	id := uuid.New()
	user := &storage.User{UserName: userName, Password: password, Uuid: id.String()}
	userStorage[userName] = user

	return id.String()
}

func (userStorage UserStorage) Get(userName string) (*storage.User, error) {
	if user, ok := userStorage[userName]; ok {
		return user, nil
	}

	return &storage.User{}, errors.New("user does not exist")
}

func NewUserStorage() UserStorageInterface {
	return &UserStorage{}
}
