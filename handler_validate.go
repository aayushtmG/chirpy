package main

import (
	"encoding/json"
	"net/http"
)


func handleChirpsValidate(w http.ResponseWriter,r *http.Request){
		type parameters struct {
			Body string `json:"body"`
		}
		type returnVals struct {
			Valid bool `json:"valid"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			respondWithError(w,http.StatusInternalServerError,"Couldn't read the parameters",err)
			return
		}
		const maxChirpLength = 140
		if len(params.Body) > maxChirpLength {
			respondWithError(w,http.StatusBadRequest,"Chirp is too long",nil)
			return
		}

		respondWithJson(w,200,returnVals{
			Valid: true,
		})

}
