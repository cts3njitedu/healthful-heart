package handlers

import (
	"net/http"
	"github.com/gorilla/context"
	"github.com/cts3njitedu/healthful-heart/models"
	"github.com/cts3njitedu/healthful-heart/security"
	"github.com/cts3njitedu/healthful-heart/utils"
)

type TokenHandler struct {
	environmentUtil utils.IEnvironmentUtility
	jwtToken security.IJwtToken
}


type ITokenHandler interface {
	GetToken(w http.ResponseWriter, r *http.Request)
}

func NewTokenHandler(environmentUtil utils.IEnvironmentUtility, jwtToken security.IJwtToken) *TokenHandler {
	return &TokenHandler{environmentUtil, jwtToken}
}

func (handler *TokenHandler) GetToken(w http.ResponseWriter, r *http.Request) {
	credentials:=context.Get(r,"credentials")
	var creds models.Credentials

	if c, ok := credentials.(models.Credentials); ok {
		creds = c
	}
	
	accessCookie, err := handler.jwtToken.CreateAccessToken(creds)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: accessCookie.Name,
		Value: accessCookie.Value,
		Expires: accessCookie.Expires,
	})
	w.Header().Set("Token", accessCookie.Value)
	refreshCookie, err := handler.jwtToken.CreateRefreshToken(creds)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: refreshCookie.Name,
		Value: refreshCookie.Value,
		Expires: refreshCookie.Expires,
		HttpOnly: refreshCookie.HttpOnly,
	})

	w.Write([]byte(`{"message": "Login Successful"}`))
}