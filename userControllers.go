package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct{
		Id uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time  `json:"updated_at"`
		Email string `json:"email"`
	}



func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter,r *http.Request) {
	type reqBody struct {
		Email string `json:"email"`
	}	
	type response struct {
		User
	}
	params := reqBody{}

	decoder := json.NewDecoder(r.Body)
	err:= decoder.Decode(&params)
	if err != nil {
		respondWithError(w,http.StatusInternalServerError,"Couldn't decode parameters",err)
		return
	}

	user,err := cfg.db.CreateUser(r.Context(),params.Email)
	if err != nil {
	respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJson(w,http.StatusCreated,response{
		User: User{
			Id:        user.ID,
			Created_at: user.CreatedAt,
			Updated_at: user.UpdatedAt,
			Email:     user.Email,
		},
	})

}

