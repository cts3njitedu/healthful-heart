package handlers

import (
	"net/http"
	"github.com/gorilla/context"
	"github.com/cts3njitedu/healthful-heart/models"
	"github.com/cts3njitedu/healthful-heart/security"
	"github.com/cts3njitedu/healthful-heart/utils"
	"strings"
	"encoding/json"
	"fmt"
)

type TokenHandler struct {
	environmentUtil utils.IEnvironmentUtility
	jwtToken security.IJwtToken
}


type ITokenHandler interface {
	GetToken(w http.ResponseWriter, r *http.Request)
}

type Exception struct {
	Message string `json:"message"`
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

func (handler *TokenHandler) ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				claims := &security.Claims{}
				credentials, err := handler.jwtToken.ValidateToken(bearerToken[1],claims)
				
				if err != nil {
					refreshToken, err := r.Cookie("refresh_token")
					fmt.Printf("Refresh Token: %+v\n", refreshToken)
					if err != nil {
						fmt.Println("Refresh token has expired...")
						w.WriteHeader(http.StatusUnauthorized)
						json.NewEncoder(w).Encode(Exception{Message: "Unauthorized access"})
						return
					}
					credentials, err = handler.jwtToken.ValidateRefreshToken(refreshToken.Value)
					if err != nil {
						w.WriteHeader(http.StatusUnauthorized)
						json.NewEncoder(w).Encode(Exception{Message: "Unauthorized access"})
						return
					}
					accessCookie, err := handler.jwtToken.CreateAccessToken(credentials)

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
					context.Set(r,"credentials", credentials)
					next.ServeHTTP(w, r)
				
				} else {
					w.Header().Set("Token", bearerToken[1])
					context.Set(r,"credentials", credentials)
					next.ServeHTTP(w, r)
				}	
			} else {
				json.NewEncoder(w).Encode(Exception{Message: "Something went wrong"})
			}

		} else {
			json.NewEncoder(w).Encode(Exception{Message: "Authorizatiton header is required"})
		}
		
		
	})
}