package main

import (
	"net/http"
	"log"
	"sync/atomic"
	"fmt"
	"os"
	"database/sql"
	"time"

	// "github.com/kitaclysm/chirpy/internal/config"
	"github.com/kitaclysm/chirpy/internal/database"
	"github.com/joho/godotenv"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits	atomic.Int32
	queries			*database.Queries
	platform		string
	jwtSecret		string
}

type LoginResponse struct {
	ID				uuid.UUID	`json:"id"`
	CreatedAt		time.Time	`json:"created_at"`
	UpdatedAt		time.Time	`json:"updated_at"`
	Email			string		`json:"email"`
	Token			string		`json:"token"`
	RefreshToken	string		`json:"refresh_token"`
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


func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("JWT_SECRET")
	
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()
	
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		queries:	dbQueries,
		platform:	platform,
		jwtSecret:		secret,
	}

	mux := http.NewServeMux()
	
	srv := &http.Server{
		Addr:		":8080",
		Handler:	mux,
	}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", handler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsGet)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChirpGet)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUserCreate)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	log.Fatal(srv.ListenAndServe())
}