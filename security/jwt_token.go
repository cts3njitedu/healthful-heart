package security

import (
	"github.com/cts3njitedu/healthful-heart/models"
	"github.com/cts3njitedu/healthful-heart/utils"
	"strconv"
	"github.com/dgrijalva/jwt-go"
	"time"
	"github.com/cts3njitedu/healthful-heart/repositories/mysqlrepo"
)
type JwtToken struct {
	environmentUtil utils.IEnvironmentUtility
	hasher IPasswordHasher
	tokenRepository mysqlrepo.ITokenRepository

}
type IJwtToken interface {
	CreateAccessToken(creds models.Credentials) (models.JwtCookie, error)
	CreateRefreshToken(creds models.Credentials) (models.JwtCookie, error)
}

type Claims struct {
	Username string `json:"username"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	UserId string `json:"userId"`
	jwt.StandardClaims
}

func NewJwtToken(environmentUtil utils.IEnvironmentUtility, hasher IPasswordHasher, tokenRepository mysqlrepo.ITokenRepository) *JwtToken {
	return &JwtToken{environmentUtil, hasher, tokenRepository}
}


func (jwtToken * JwtToken) CreateAccessToken(creds models.Credentials) (models.JwtCookie, error) {
	accessTokenExpireTime,err:=strconv.Atoi(jwtToken.environmentUtil.GetEnvironmentString("JWT_ACCESS_TOKEN_EXPIRATION_TIME"))
	
	if err!=nil {
		panic(err.Error())
	}

	expirationTime := time.Now().Add(time.Duration(accessTokenExpireTime) * time.Minute)
	

	claims := &Claims{
		Username: creds.Username,
		FirstName: creds.FirstName,
		LastName: creds.LastName,
		UserId: creds.UserId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenSecret:=jwtToken.environmentUtil.GetEnvironmentString("JWT_SECRET_KEY")
	tokenString, err := token.SignedString([]byte(tokenSecret))
	
	var cookie models.JwtCookie
	if err != nil {
		return models.JwtCookie{}, err
	}

	cookie.Name = "access_token"
	cookie.Value = tokenString
	cookie.Expires = expirationTime

	return cookie, nil

}

func (jwtToken * JwtToken) CreateRefreshToken(creds models.Credentials) (models.JwtCookie, error) {

	refreshTokenExpireTime,err:=strconv.Atoi(jwtToken.environmentUtil.GetEnvironmentString("JWT_REFRESH_TOKEN_EXPIRATION_TIME"))

	if err!=nil {
		panic(err.Error())
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256);
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	expirationTime := time.Now().Add(time.Hour * time.Duration(refreshTokenExpireTime))
	tokenSecret:=jwtToken.environmentUtil.GetEnvironmentString("JWT_SECRET_KEY")
	rtClaims["sub"] = creds.UserId
	rtClaims["exp"] = expirationTime.Unix()
	refreshTokenString, err := refreshToken.SignedString([]byte(tokenSecret))

	if err != nil {
		return models.JwtCookie{}, err
	}

	var cookie models.JwtCookie
	cookie.Name = "refresh_token"
	cookie.Value = refreshTokenString
	cookie.Expires = expirationTime
	cookie.HttpOnly = true

	hashedToken, err := jwtToken.hasher.HashPassword(refreshTokenString)

	if err != nil {
		return models.JwtCookie{}, err
	}

	err = jwtToken.tokenRepository.SaveRefreshToken(hashedToken, expirationTime, creds.UserId)

	if err != nil {
		return models.JwtCookie{}, err
	}


	return cookie, nil

}