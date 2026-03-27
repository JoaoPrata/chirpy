package main

import (
	_ "github.com/lib/pq"
	"log"
	"os"
	"net/http"
	"sync/atomic"
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/JoaoPrata/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries *database.Queries
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DB_URL")
	isDev := os.Getenv("PLATFORM") == "dev"
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error opening db")
	}

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries: database.New(db),
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	if isDev {
		mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	}else {
		mux.HandleFunc("POST /admin/reset", apiCfg.handlerDisabledReset)
	}
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.handlerValidate)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsers)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}