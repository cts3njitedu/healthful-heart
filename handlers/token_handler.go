package handlers

import (
	"net/http"
	"github.com/gorilla/context"
	"github.com/dgrijalva/jwt-go"
	"github.com/cts3njitedu/healthful-heart/models"
	"time"
)

type TokenHandler struct {}


type ITokenHandler interface {
	GetToken(w http.ResponseWriter, r *http.Request)
}

func NewTokenHandler() *TokenHandler {
	return &TokenHandler{}
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
func (handler *TokenHandler) GetToken(w http.ResponseWriter, r *http.Request) {
	credentials:=context.Get(r,"credentials")
	var creds models.Credentials

	if c, ok := credentials.(models.Credentials); ok {
		creds = c
	}
	

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}


	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name: "access_token",
		Value: tokenString,
		Expires: expirationTime,
	})
	
	refreshToken := jwt.New(jwt.SigningMethodHS256);
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshExpirationTime := time.Now().Add(time.Hour * 24)
	rtClaims["sub"] = 1
	rtClaims["exp"] = refreshExpirationTime.Unix()
	refreshTokenString, err := refreshToken.SignedString([]byte("secret"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: "refresh_token",
		Value: refreshTokenString,
		Expires: refreshExpirationTime,
		HttpOnly: true,
	})
	w.Write([]byte(`{"message": "Login Successful"}`))
}