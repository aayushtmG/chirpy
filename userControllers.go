package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)



func (cfg apiConfig) handleCreateUser(w http.ResponseWriter,r *http.Request) {
	type reqBody struct {
		Email string `json:"email"`
	}	
	params := reqBody{}

	data,err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Error reading body: %s\n",err)
	}

	err = json.Unmarshal(data,&params)
	if err != nil {
		log.Fatalf("Error unmarshalling body: %s\n",err)
	}

	user, err := cfg.db.CreateUser(r.Context(),params.Email)
	if err != nil {
		log.Fatalf("Error creating user: %s\n",err)
	}
	userData := User{
		Id: user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email: user.Email,
	}
	
	userJson,err := json.Marshal(userData)

	if err != nil {
		log.Fatalf("Error marshalling user data: %s\n",err)
	}

	w.WriteHeader(http.StatusCreated)	
	w.Write(userJson)
}

