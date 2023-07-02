package userstorage_test

import (
	"httpserver/internal/storage"
	"httpserver/internal/storage/userstorage"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserStorage_Add(t *testing.T) {
	storage := userstorage.NewUserStorage()

	userID := storage.Add("JohnDoe", "password123")

	assert.NotEmpty(t, userID)
	user, _ := storage.Get("JohnDoe")
	assert.Equal(t, "JohnDoe", user.UserName)
	assert.Equal(t, "password123", user.Password)
}

func TestUserStorage_Get(t *testing.T) {
	storage := userstorage.NewUserStorage()

	storage.Add("JohnDoe", "password123")

	user, err := storage.Get("JohnDoe")
	assert.NoError(t, err)
	assert.Equal(t, "JohnDoe", user.UserName)
	assert.Equal(t, "password123", user.Password)
}

func TestUserStorage_Get_NonExistentUser(t *testing.T) {
	storageInstance := userstorage.NewUserStorage()

	user, err := storageInstance.Get("NonExistentUser")
	assert.Error(t, err)
	assert.Equal(t, &storage.User{}, user)
}

func BenchmarkAdd(b *testing.B) {
	storage := userstorage.NewUserStorage()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		storage.Add("john.doe", "password")
	}
}

func BenchmarkGet(b *testing.B) {
	storage := userstorage.NewUserStorage()
	storage.Add("john.doe", "password")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		storage.Get("john.doe")
	}
}
