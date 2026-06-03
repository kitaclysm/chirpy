package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "some-secret"
	expiresIn := time.Hour

	tokenString, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("unexpected error making JWT: %v", err)
	}

	gotUserID, err := ValidateJWT(tokenString, secret)
	if err != nil {
		t.Fatalf("unexpected error validating JWT: %v", err)
	}

	if gotUserID != userID {
		t.Fatalf("expected user ID %v, got %v", userID, gotUserID)
	}
}

func TestMakeAndValidateExpiredJWT(t *testing.T) {
	userID := uuid.New()
	secret := "some-secret"
	expiresIn := -time.Hour

	tokenString, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("unexpected error making JWT: %v", err)
	}

	_, err = ValidateJWT(tokenString, secret)
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}

func TestWrongSecretJWT(t *testing.T) {
	userID := uuid.New()
	secret := "some-secret"
	expiresIn := time.Hour

	tokenString, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("unexpected error making JWT: %v", err)
	}

	_, err = ValidateJWT(tokenString, "another-secret")
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}