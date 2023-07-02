package tokenstorage

import (
	"errors"
	"httpserver/internal/storage"
)

type TokenStorageInterface interface {
	Add(string, *storage.User)
	Get(string) (*storage.User, error)
	Delete(token string)
}

type TokenStorage map[string]*storage.User

func (tokenStorage TokenStorage) Add(token string, user *storage.User) {
	tokenStorage[token] = user
}

func (tokenStorage TokenStorage) Get(token string) (*storage.User, error) {
	if userData, ok := tokenStorage[token]; ok {
		delete(tokenStorage, token)
		return userData, nil
	}

	return &storage.User{}, errors.New("user does not exist")
}

func (tokenStorage TokenStorage) Delete(token string) {
	delete(tokenStorage, token)
}

func NewTokenStorage() TokenStorageInterface {
	return &TokenStorage{}
}
