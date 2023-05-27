package userstorage

import (
	"errors"

	"github.com/google/uuid"
)

type User struct {
	UserName string
	Password string
	Uuid     string
}

type UserStorage map[string]*User

func (userStorage UserStorage) Add(userName string, password string) string {
	id := uuid.New()
	user := &User{UserName: userName, Password: password, Uuid: id.String()}
	userStorage[userName] = user

	return id.String()
}

func (userStorage UserStorage) Get(userName string) (*User, error) {
	if userData, ok := userStorage[userName]; ok {
		return userData, nil
	}

	return &User{}, errors.New("user does not exist")
}

func NewUserStorage() *UserStorage {
	return &UserStorage{}
}

type ActiveUsersStorage map[string]*User

func (activeUsersStorage ActiveUsersStorage) Add(user *User) {
	activeUsersStorage[user.UserName] = user
}

func (activeUsersStorage ActiveUsersStorage) Get(userName string) (*User, error) {
	if userData, ok := activeUsersStorage[userName]; ok {
		return userData, nil
	}

	return &User{}, errors.New("user does not exist")
}

func (activeUsersStorage ActiveUsersStorage) Delete(user *User) {
	delete(activeUsersStorage, user.UserName)
}

func NewActiveUsersStorage() *ActiveUsersStorage {
	return &ActiveUsersStorage{}
}
