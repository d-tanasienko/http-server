package activeuserstorage

import (
	"errors"
	"httpserver/internal/storage"
)

type ActiveUsersStorageInterface interface {
	Add(*storage.User)
	Get(string) (*storage.User, error)
	Delete(user *storage.User)
	GetNames() []string
}

type ActiveUsersStorage map[string]*storage.User

func (activeUsersStorage ActiveUsersStorage) Add(user *storage.User) {
	activeUsersStorage[user.UserName] = user
}

func (activeUsersStorage ActiveUsersStorage) Get(userName string) (*storage.User, error) {
	if userData, ok := activeUsersStorage[userName]; ok {
		return userData, nil
	}

	return &storage.User{}, errors.New("user does not exist")
}

func (activeUsersStorage ActiveUsersStorage) Delete(user *storage.User) {
	delete(activeUsersStorage, user.UserName)
}

func (activeUsersStorage ActiveUsersStorage) GetNames() []string {
	userNames := make([]string, len(activeUsersStorage))

	i := 0
	for k := range activeUsersStorage {
		userNames[i] = k
		i++
	}

	return userNames
}

func NewActiveUsersStorage() ActiveUsersStorageInterface {
	return &ActiveUsersStorage{}
}
