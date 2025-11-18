package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/aayushtmG/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	platform string	
}

func (cfg *apiConfig) middlewareMetricInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		cfg.fileserverHits.Add(1)
		w.Header().Set("Cache-Control","no-cache")	
		next.ServeHTTP(w,r)
	})
}




func main() {
	godotenv.Load()
	const filepathRoot = "."
	const port = "8080"
	dbUrl := os.Getenv("DB_URL")
	platform := os.Getenv("platform")
	db,err := sql.Open("postgres",dbUrl)
	if err != nil {
		log.Fatal("Error connecting database")
	}
	dbQueries  := database.New(db)

	apiCfg := apiConfig{
		db: dbQueries,
		platform: platform,
	}

	mux :=  http.NewServeMux()

	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	
	mux.Handle("/app/",apiCfg.middlewareMetricInc(
		http.StripPrefix("/app",
		http.FileServer(
		http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz/",handlerReadiness)
	mux.HandleFunc("GET /admin/metrics",apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset",apiCfg.handlerReset)
	mux.HandleFunc("POST /api/users",apiCfg.handleCreateUser)
	mux.HandleFunc("POST /api/chirps",apiCfg.handleCreateChirps)
	mux.HandleFunc("GET /api/chirps",apiCfg.handleGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpId}",apiCfg.handleGetOneChirps)

	log.Printf("Serving on port: %s\n",port)
	log.Fatal(server.ListenAndServe())
}


