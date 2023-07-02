package tokenstorage_test

import (
	"httpserver/internal/storage"
	"httpserver/internal/storage/tokenstorage"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenStorage_Add(t *testing.T) {
	tokenStorageInstance := tokenstorage.NewTokenStorage()

	user := &storage.User{UserName: "JohnDoe", Password: "password123"}
	token := "abc123"

	tokenStorageInstance.Add(token, user)

	userInStorage, _ := tokenStorageInstance.Get(token)
	assert.Equal(t, user, userInStorage)
}

func TestTokenStorage_Get(t *testing.T) {
	tokenStorageInstance := tokenstorage.NewTokenStorage()

	user := &storage.User{UserName: "JohnDoe", Password: "password123"}
	token := "abc123"

	tokenStorageInstance.Add(token, user)

	resultUser, err := tokenStorageInstance.Get(token)

	assert.NoError(t, err)
	assert.Equal(t, user, resultUser)

	_, err = tokenStorageInstance.Get(token)

	assert.Error(t, err)
}

func TestTokenStorage_Delete(t *testing.T) {
	tokenStorageInstance := tokenstorage.NewTokenStorage()

	user := &storage.User{UserName: "JohnDoe", Password: "password123"}
	token := "abc123"

	tokenStorageInstance.Add(token, user)

	tokenStorageInstance.Delete(token)

	_, err := tokenStorageInstance.Get(token)

	assert.Error(t, err)
}

func TestTokenStorage_Get_NonExistentToken(t *testing.T) {
	tokenStorageInstance := tokenstorage.NewTokenStorage()

	user, err := tokenStorageInstance.Get("nonexistent-token")

	assert.Error(t, err)
	assert.Equal(t, &storage.User{}, user)
}

func BenchmarkAdd(b *testing.B) {
	tokenStorageInstance := tokenstorage.NewTokenStorage()
	user := &storage.User{UserName: "John Doe"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tokenStorageInstance.Add("token", user)
	}
}

func BenchmarkGet(b *testing.B) {
	tokenStorageInstance := tokenstorage.NewTokenStorage()
	user := &storage.User{UserName: "John Doe"}
	tokenStorageInstance.Add("token", user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tokenStorageInstance.Get("token")
	}
}

func BenchmarkDelete(b *testing.B) {
	tokenStorageInstance := tokenstorage.NewTokenStorage()
	user := &storage.User{UserName: "John Doe"}
	tokenStorageInstance.Add("token", user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tokenStorageInstance.Delete("token")
	}
}
