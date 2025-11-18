package main

import (
	"encoding/json"
	// "fmt"
	"net/http"
)


func handleChirpsValidate(w http.ResponseWriter,r *http.Request){
		type parameters struct {
			Body string `json:"body"`
		}
		type returnVals struct {
			CleanedBody string `json:"cleaned_body"`
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
		badWords := map[string]struct{}{
			"kerfuffle": {},
			"sharbert": {},
			"fornax": {},
		}
		cleaned := getCleanedBody(params.Body,badWords)
		respondWithJson(w,200,returnVals{
			CleanedBody: cleaned,
		})

}

//this was my solution
// func filterBadWord(s string) string{
// 	splited := strings.Split(s," ")
// 	for in, str := range splited{
// 		l := strings.ToLower(str)
// 		if l == "kerfuffle" || l == "sharbert" ||  l == "fornax" {
// 			splited[in] = "****"
// 		}
// 	}
// 	return strings.Join(splited, " ")
// }
// func getCleanedBody(body string, badWords map[string]struct{}) string {
// 	words := strings.Split(body, " ")
// 	for i, word := range words {
// 		loweredWord := strings.ToLower(word)
// 		if _, ok := badWords[loweredWord]; ok {
// 			words[i] = "****"
// 		}
// 	}
// 	cleaned := strings.Join(words, " ")
// 	return cleaned
// }
