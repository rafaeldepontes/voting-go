package service

import (
	"testing"
)

func TestHashAndVerifyPassword(t *testing.T) {
	password := "securePassword123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if hash == "" {
		t.Fatal("hash should not be empty")
	}

	match, err := VerifyPassword(password, hash)
	if err != nil {
		t.Fatalf("failed to verify password: %v", err)
	}
	if !match {
		t.Fatal("password should match")
	}

	match, err = VerifyPassword("wrongPassword", hash)
	if err != nil {
		t.Fatalf("failed to verify password: %v", err)
	}
	if match {
		t.Fatal("wrong password should not match")
	}
}

func TestVerifyPasswordInvalidFormat(t *testing.T) {
	_, err := VerifyPassword("password", "invalidHash")
	if err == nil {
		t.Fatal("should fail with invalid hash format")
	}
}
