package activeuserstorage_test

import (
	"fmt"
	"httpserver/internal/storage"
	"httpserver/internal/storage/activeuserstorage"
)

func ExampleActiveUsersStorage_Add() {
	activeUsersStorage := activeuserstorage.NewActiveUsersStorage()
	user := &storage.User{UserName: "john.doe", Password: "password"}

	activeUsersStorage.Add(user)
	user, _ = activeUsersStorage.Get("john.doe")
	fmt.Println(user)

	// Output: &{john.doe password }
}

func ExampleActiveUsersStorage_Get() {
	activeUsersStorage := activeuserstorage.NewActiveUsersStorage()
	user := &storage.User{UserName: "john.doe", Password: "password"}
	activeUsersStorage.Add(user)

	user, _ = activeUsersStorage.Get("john.doe")
	fmt.Println(user)

	// Output: &{john.doe password }
}

func ExampleActiveUsersStorage_Delete() {
	activeUsersStorage := activeuserstorage.NewActiveUsersStorage()
	user := &storage.User{UserName: "john.doe", Password: "password"}
	activeUsersStorage.Add(user)

	activeUsersStorage.Delete(user)
	fmt.Println(activeUsersStorage.Get("john.doe"))

	// Output: &{  } user does not exist
}

func ExampleActiveUsersStorage_GetNames() {
	activeUsersStorage := activeuserstorage.NewActiveUsersStorage()
	user1 := &storage.User{UserName: "john.doe", Password: "password"}
	user2 := &storage.User{UserName: "jane.doe", Password: "password"}
	activeUsersStorage.Add(user1)
	activeUsersStorage.Add(user2)

	fmt.Println(activeUsersStorage.GetNames())

	// Output: [john.doe jane.doe]
}
