package tokenstorage_test

import (
	"fmt"
	"httpserver/internal/storage"
	"httpserver/internal/storage/tokenstorage"
)

func ExampleTokenStorage_Add() {
	tokenStorage := tokenstorage.NewTokenStorage()
	user := &storage.User{UserName: "john.doe", Password: "password"}

	tokenStorage.Add("token123", user)
	user, _ = tokenStorage.Get("token123")
	fmt.Println(user)

	// Output: &{john.doe password }
}

func ExampleTokenStorage_Get() {
	tokenStorage := tokenstorage.NewTokenStorage()
	user := &storage.User{UserName: "john.doe", Password: "password"}
	tokenStorage.Add("token123", user)

	user, _ = tokenStorage.Get("token123")
	fmt.Println(user)

	// Output: &{john.doe password }
}

func ExampleTokenStorage_Delete() {
	tokenStorage := tokenstorage.NewTokenStorage()
	user := &storage.User{UserName: "john.doe", Password: "password"}
	tokenStorage.Add("token123", user)

	tokenStorage.Delete("token123")

	fmt.Println(tokenStorage.Get("token123"))

	// Output: &{  } user does not exist
}
