package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		cfg.fileserverHits.Add(1)
		w.Header().Set("Cache-Control","no-cache")	
		next.ServeHTTP(w,r)
	})
}



func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{}

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
	mux.HandleFunc("POST /api/validate_chirp",apiCfg.handlerValidateChirp)

	log.Printf("Serving on port: %s\n",port)
	log.Fatal(server.ListenAndServe())
}



func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter,r *http.Request){
	w.Header().Add("Content-Type","application/json")
	
	 params := struct{
		Body string `json:"body"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		errResp := struct{
			Error string `json:"error"`
		}{
			Error: "Something went wrong",
		}
		data, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(data)
		return
	}
	if len(params.Body) > 140 {
		errResp := struct{
			Error string `json:"error"`
		}{
			Error: "Chirp is too long",
		}
		data, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(data)
		return
	}



		resp := struct{
			Valid bool `json:"valid"`
		}{
			Valid: true,
		}
		data, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusOK)
		w.Write(data)
}


func handlerReadiness(w http.ResponseWriter,r *http.Request){
	w.Header().Add("Content-Type","text/plain; charset=utf-8")	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}


func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`<html>
	<body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
	</body>
	</html>`,cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter,r *http.Request){
	cfg.fileserverHits.Store(0)		
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
