package main

import (
	"database/sql"

	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/faroukhamissou-dev/chirpy/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	secretKey := os.Getenv("SECRET_KEY")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf(err.Error())
	}
	dbQueries := database.New(db)
	platform := os.Getenv("PLATFORM")
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileServerHits: atomic.Int32{},
		dbQueries:      dbQueries,
		PLATFORM:       platform,
		SECRET_KEY: secretKey,
	}

	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)
	mux.HandleFunc("GET /api/healthz", healthz)
	mux.HandleFunc("GET /admin/metrics", apiCfg.countHits)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHits)
	mux.HandleFunc("POST /api/chirps", apiCfg.addChirp)
	mux.HandleFunc("POST /api/users", apiCfg.addUser)
	mux.HandleFunc("GET /api/chirps", apiCfg.getChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirp)
	mux.HandleFunc("POST /api/login", apiCfg.login)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())

}
