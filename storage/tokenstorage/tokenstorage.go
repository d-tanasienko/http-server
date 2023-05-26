package tokenstorage

import (
	"errors"
	"httpserver/storage/userstorage"
)

type TokenStorage map[string]*userstorage.User

func (tokenStorage TokenStorage) Add(token string, user *userstorage.User) {
	tokenStorage[token] = user
}

func (tokenStorage TokenStorage) Get(token string) (*userstorage.User, error) {
	if userData, ok := tokenStorage[token]; ok {
		delete(tokenStorage, token)
		return userData, nil
	}

	return &userstorage.User{}, errors.New("user does not exist")
}

func (tokenStorage TokenStorage) Delete(token string) {
	delete(tokenStorage, token)
}

func NewTokenStorage() *TokenStorage {
	return &TokenStorage{}
}
