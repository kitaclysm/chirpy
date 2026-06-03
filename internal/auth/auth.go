package auth

import (
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	userIDString := userID.String()

	issTime := time.Now().UTC()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:	"chirpy-access",
		IssuedAt: jwt.NewNumericDate(issTime),
		ExpiresAt: jwt.NewNumericDate(issTime.Add(expiresIn)),
		Subject: userIDString,
	})
	secretBytes := []byte(tokenSecret)
	signedToken, err := token.SignedString(secretBytes)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}