package handlers

import (
	"net/http"
	"github.com/gorilla/context"
	"github.com/dgrijalva/jwt-go"
	"github.com/cts3njitedu/healthful-heart/models"
	"time"
	"github.com/cts3njitedu/healthful-heart/utils"
	"strconv"
)

type TokenHandler struct {
	environmentUtil utils.IEnvironmentUtility
}


type ITokenHandler interface {
	GetToken(w http.ResponseWriter, r *http.Request)
}

func NewTokenHandler(environmentUtil utils.IEnvironmentUtility) *TokenHandler {
	return &TokenHandler{environmentUtil}
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
	
	accessTokenExpireTime,err:=strconv.Atoi(handler.environmentUtil.GetEnvironmentString("JWT_ACCESS_TOKEN_EXPIRATION_TIME"))
	
	if err!=nil {
		panic(err.Error())
	}

	refreshTokenExpireTime,err:=strconv.Atoi(handler.environmentUtil.GetEnvironmentString("JWT_REFRESH_TOKEN_EXPIRATION_TIME"))

	if err!=nil {
		panic(err.Error())
	}
	
	expirationTime := time.Now().Add(time.Duration(accessTokenExpireTime) * time.Minute)
	

	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}


	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenSecret:=handler.environmentUtil.GetEnvironmentString("JWT_SECRET_KEY")
	tokenString, err := token.SignedString([]byte(tokenSecret))
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
	refreshExpirationTime := time.Now().Add(time.Hour * time.Duration(refreshTokenExpireTime))
	rtClaims["sub"] = 1
	rtClaims["exp"] = refreshExpirationTime.Unix()
	refreshTokenString, err := refreshToken.SignedString([]byte(tokenSecret))
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