package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	res := Response{
		Msg:  "Server is running",
		Code: 200,
	}

	jsonResponse, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
