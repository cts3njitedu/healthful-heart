package handlers

import (
	"net/http"
	"encoding/json"
	"github.com/cts3njitedu/healthful-heart/repositories/mongorepo"
	"github.com/cts3njitedu/healthful-heart/models"
	"github.com/cts3njitedu/healthful-heart/utils"
	"fmt"
)

type LoginResponse struct {
	IsError bool `json:"isError"`
	Message string `json:"Message"`
}




func GetLoginPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type","application/json")
	page:=mongorepo.GetPage("LOGIN_PAGE")
	json.NewEncoder(w).Encode(page)
}


func PostLoginPage(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials;
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
	}

	hashed, err := utils.HashPassword(credentials.Password);

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
        return
	}
	
	fmt.Printf("Password: %s Hashed Password: %s", credentials.Password, hashed)
	w.Write([]byte(`{"message": "Successfully Created"}`))
}