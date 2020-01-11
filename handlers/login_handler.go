package handlers

import (
	"net/http"
	"encoding/json"
	"github.com/cts3njitedu/healthful-heart/models"
	"github.com/cts3njitedu/healthful-heart/services"
)

type LoginResponse struct {
	IsError bool `json:"isError"`
	Message string `json:"Message"`
}

type LoginHandler struct {
	authenticationService services.IAuthenticationService
}

type ILoginHandler interface {
	GetLoginPage(w http.ResponseWriter, r *http.Request)
	PostLoginPage(w http.ResponseWriter, r *http.Request)
}
func NewLoginHandler(authenticationService services.IAuthenticationService) *LoginHandler {
	return &LoginHandler{authenticationService}
}

func (handler *LoginHandler) GetLoginPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type","application/json")
	page:=handler.authenticationService.GetAuthenticationPage("LOGIN");
	json.NewEncoder(w).Encode(page)
}


func (handler *LoginHandler) PostLoginPage(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials;
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
	}
	
	w.Write([]byte(`{"message": "Successfully Created"}`))
}