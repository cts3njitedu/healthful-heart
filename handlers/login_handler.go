package handlers

import (
	"net/http"
	"encoding/json"
	
)

type LoginResponse struct {
	IsError bool `json:"isError"`
	Message string `json:"Message"`
}


func Login(w http.ResponseWriter, r *http.Request) {

	loginResponse:= LoginResponse {
		IsError: false,
		Message: "Login Successful",
	}
	json.NewEncoder(w).Encode(loginResponse)
}