package activeuserstorage_test

import (
	"httpserver/internal/storage"
	"httpserver/internal/storage/activeuserstorage"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActiveUsersStorage_Add(t *testing.T) {
	activeUsersStorage := activeuserstorage.NewActiveUsersStorage()

	user := &storage.User{UserName: "JohnDoe", Password: "password123"}

	activeUsersStorage.Add(user)
	userInStorage, _ := activeUsersStorage.Get("JohnDoe")
	assert.Equal(t, user, userInStorage)
}

func TestActiveUsersStorage_Get(t *testing.T) {
	activeUsersStorage := activeuserstorage.NewActiveUsersStorage()

	user := &storage.User{UserName: "JohnDoe", Password: "password123"}

	activeUsersStorage.Add(user)

	resultUser, err := activeUsersStorage.Get("JohnDoe")

	assert.NoError(t, err)
	assert.Equal(t, user, resultUser)
}

func TestActiveUsersStorage_Get_NonExistentUser(t *testing.T) {
	activeUsersStorage := activeuserstorage.NewActiveUsersStorage()

	user, err := activeUsersStorage.Get("NonExistentUser")

	assert.Error(t, err)
	assert.Equal(t, &storage.User{}, user)
}

func TestActiveUsersStorage_Delete(t *testing.T) {
	activeUsersStorage := activeuserstorage.NewActiveUsersStorage()

	user := &storage.User{UserName: "JohnDoe", Password: "password123"}

	activeUsersStorage.Add(user)
	activeUsersStorage.Delete(user)

	_, err := activeUsersStorage.Get("JohnDoe")

	assert.Error(t, err)
}

func TestActiveUsersStorage_GetNames(t *testing.T) {
	activeUsersStorage := activeuserstorage.NewActiveUsersStorage()

	user1 := &storage.User{UserName: "JohnDoe", Password: "password123"}
	user2 := &storage.User{UserName: "JaneSmith", Password: "password456"}

	activeUsersStorage.Add(user1)
	activeUsersStorage.Add(user2)

	expectedNames := []string{"JohnDoe", "JaneSmith"}
	userNames := activeUsersStorage.GetNames()

	assert.ElementsMatch(t, expectedNames, userNames)
}

func BenchmarkAdd(b *testing.B) {
	activeUserStorage := activeuserstorage.NewActiveUsersStorage()
	user := &storage.User{
		UserName: "testuser",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		activeUserStorage.Add(user)
	}
}

func BenchmarkGet(b *testing.B) {
	activeUserStorage := activeuserstorage.NewActiveUsersStorage()
	user := &storage.User{
		UserName: "testuser",
	}
	activeUserStorage.Add(user)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := activeUserStorage.Get("testuser")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDelete(b *testing.B) {
	activeUserStorage := activeuserstorage.NewActiveUsersStorage()
	user := &storage.User{
		UserName: "testuser",
	}
	activeUserStorage.Add(user)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		activeUserStorage.Delete(user)
	}
}

func BenchmarkGetNames(b *testing.B) {
	activeUserStorage := activeuserstorage.NewActiveUsersStorage()
	user := &storage.User{
		UserName: "testuser",
		// Set other user properties if needed
	}
	activeUserStorage.Add(user)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = activeUserStorage.GetNames()
	}
}
