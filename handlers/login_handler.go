package handlers

import (
	"net/http"
	"encoding/json"
	"github.com/cts3njitedu/healthful-heart/repositories/mongorepo"
)

type LoginResponse struct {
	IsError bool `json:"isError"`
	Message string `json:"Message"`
}




func GetLoginPage(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type","application/json")
	page:=mongorepo.GetPage("LOGIN_PAGE")
	json.NewEncoder(response).Encode(page)
	
}