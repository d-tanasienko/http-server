package userstorage_test

import (
	"httpserver/internal/storage/userstorage"
)

func ExampleUserStorage_Add() {
	storage := userstorage.NewUserStorage()

	storage.Add("john.doe", "password")

	user, _ := storage.Get("john.doe")
	_ = user
}

func ExampleUserStorage_Get() {
	storage := userstorage.NewUserStorage()
	storage.Add("john.doe", "password")

	user, _ := storage.Get("john.doe")
	_ = user
}
