package main

import (
	"time"
	"sync/atomic"

	"github.com/kitaclysm/chirpy/internal/database"
	"github.com/google/uuid"
)

type apiConfig struct {
	fileserverHits	atomic.Int32
	queries			*database.Queries
	platform		string
	jwtSecret		string
}

type ReturnUser struct {
	ID				uuid.UUID	`json:"id"`
	CreatedAt		time.Time	`json:"created_at"`
	UpdatedAt		time.Time	`json:"updated_at"`
	Email			string		`json:"email"`
	Token			string		`json:"token"`
	RefreshToken	string		`json:"refresh_token"`
	IsChirpyRed		bool		`json:"is_chirpy_red"`
}

type Parameters struct {
	Email 				string `json:"email"`
	Password 			string `json:"password"`
	// ExpiresInSeconds 	int `json:"expires_in_seconds"`
}

type returnChirp struct {
	ID			uuid.UUID	`json:"id"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
	Body		string		`json:"body"`
	UserID		uuid.UUID	`json:"user_id"`
}

type upgradeRequest struct {
	Event		string `json:"event"`
	Data		struct	{
		UserID	uuid.UUID `json:"user_id"`
	} `json:"data"`
}