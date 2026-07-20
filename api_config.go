package main

import (
	"sync/atomic"

	"github.com/faroukhamissou-dev/chirpy/internal/database"
)
type apiConfig struct {
	fileServerHits atomic.Int32
	dbQueries      *database.Queries
	PLATFORM string
	SECRET_KEY string
}
