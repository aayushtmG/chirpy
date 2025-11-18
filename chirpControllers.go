package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/aayushtmG/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct{
		Id uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time  `json:"updated_at"`
		UserId uuid.UUID `json:"user_id"`
		Body string `json:"body"`
	}

func (cfg *apiConfig) handleCreateChirps(w http.ResponseWriter,r *http.Request){
	type reqBody struct {
		Body string `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}	
	type response struct {
		Chirp	
	}

	decoder := json.NewDecoder(r.Body)
	params := reqBody{}
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

		badWords := map[string]struct{}{
			"kerfuffle": {},
			"sharbert": {},
			"fornax": {},
		}
		cleanedBody := getCleanedBody(params.Body,badWords)

		chirp, err := cfg.db.CreateChirp(r.Context(),database.CreateChirpParams{
			Body: cleanedBody,
			UserID: params.UserId,
		})		
		if err != nil {
	respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
	return
		}


		respondWithJson(w,http.StatusCreated,response{
			Chirp: Chirp{
				Id: chirp.ID,
				Created_at: chirp.CreatedAt,
				Updated_at: chirp.UpdatedAt,
				UserId: chirp.UserID,
				Body: chirp.Body,
			},
		})

}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}


func (cfg *apiConfig) handleGetAllChirps(w http.ResponseWriter,r *http.Request){

	data, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w,http.StatusInternalServerError,"Error fetching chirps",err)
		return
	}
	chirps := make([]Chirp,len(data))
	for _, dat := range data{
		chirps = append(chirps,Chirp{
			Id: dat.ID,
			Created_at: dat.CreatedAt,
			Updated_at: dat.UpdatedAt,
			UserId: dat.UserID,
			Body: dat.Body,	
		})
	}
	respondWithJson(w,200,chirps)		
}


func (cfg *apiConfig) handleGetOneChirps(w http.ResponseWriter,r *http.Request){
	chirpIdParsed,err := uuid.Parse(r.PathValue("chirpId"))
	if err != nil {
		respondWithError(w,http.StatusInternalServerError,"Error parsing params to uuid",err)
		return
	}
	dat,err := cfg.db.GetOneChirp(r.Context(),chirpIdParsed)
	if err != nil {
		respondWithError(w,http.StatusNotFound,"Chirp not found",err)
		return
	}

	respondWithJson(w,http.StatusOK,Chirp{
			Id: dat.ID,
			Created_at: dat.CreatedAt,
			Updated_at: dat.UpdatedAt,
			UserId: dat.UserID,
			Body: dat.Body,	
	})
}