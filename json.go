package main

import (
	"encoding/json"
	"log"
	"net/http"
)



func respondWithError(w http.ResponseWriter, code int,msg string,err error){
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5xx error: %s",msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJson(w,code,errorResponse{
		Error: msg,
	})
}

func respondWithJson(w http.ResponseWriter,code int,payload interface{}){
	w.Header().Set("Content-Type","application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling JSON: %s",err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(code)
	w.Write(data)
}