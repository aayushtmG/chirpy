package main

import (
	"log"
	"net/http"
)

const DEV_PLATFORM  = "dev"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter,r *http.Request){
	cfg.fileserverHits.Store(0)		

	if cfg.PLATFORM != DEV_PLATFORM {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(http.StatusText(http.StatusForbidden)))
		return
	}

	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		log.Fatalf("Error deleting users: %s\n",err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}