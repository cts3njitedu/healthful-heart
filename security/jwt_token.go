package security

import (
	"github.com/cts3njitedu/healthful-heart/models"
	"github.com/cts3njitedu/healthful-heart/mappers"
	"github.com/cts3njitedu/healthful-heart/utils"
	"strconv"
	"github.com/dgrijalva/jwt-go"
	"time"
	"github.com/cts3njitedu/healthful-heart/repositories/mysqlrepo"
	"errors"
	"fmt"
)
type JwtToken struct {
	environmentUtil utils.IEnvironmentUtility
	hasher IPasswordHasher
	tokenRepository mysqlrepo.ITokenRepository
	userRepository mysqlrepo.IUserRepository
	mapperUtil mappers.IMapper

}
type IJwtToken interface {
	CreateAccessToken(creds models.Credentials) (models.JwtCookie, error)
	CreateRefreshToken(creds models.Credentials) (models.JwtCookie, error)
	ValidateToken(bearerToken string, claims *Claims) (models.Credentials, error)
	ValidateRefreshToken(bearerToken string) (models.Credentials, error)
}

type Claims struct {
	Username string `json:"username"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	UserId string `json:"userId"`
	jwt.StandardClaims
}

func NewJwtToken(environmentUtil utils.IEnvironmentUtility, hasher IPasswordHasher, 
	tokenRepository mysqlrepo.ITokenRepository, 
	userRepository mysqlrepo.IUserRepository, mapperUtil mappers.IMapper) *JwtToken {
	return &JwtToken{environmentUtil, hasher, tokenRepository, userRepository, mapperUtil}
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

func (jwtToken * JwtToken) ValidateToken(bearerToken string, claims *Claims) (models.Credentials, error) {
	token, err := jwt.ParseWithClaims(bearerToken, claims, func(token *jwt.Token) (interface{}, error){
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("There was an error")
		}
		tokenSecret:=jwtToken.environmentUtil.GetEnvironmentString("JWT_SECRET_KEY")
		return []byte(tokenSecret), nil
	})

	if err != nil {
		if err.(*jwt.ValidationError).Errors&jwt.ValidationErrorExpired != 0 {
			fmt.Println("There is an expired token")
		}
		return models.Credentials{}, err
	}
	
	if token.Valid {
		credentials := models.Credentials {
			Username: claims.Username,
			FirstName: claims.FirstName,
			LastName: claims.LastName,
			UserId: claims.UserId,
		}
		return credentials, nil
	} else{
		return models.Credentials{}, errors.New("Unauthorized access")
	}
}

func (jwtToken * JwtToken) ValidateRefreshToken(bearerToken string) (models.Credentials, error) {
	claims := &Claims{}
	_, err := jwtToken.ValidateToken(bearerToken, claims)
	if err != nil {
		return models.Credentials{}, err
	}
	userId := claims.Subject;

	userToken, err := jwtToken.userRepository.GetUserToken(userId)

	if err != nil {
		return models.Credentials{}, err
	}
	
	err = jwtToken.hasher.CompareHashWithPassword(userToken.RefreshToken, bearerToken)

	if err != nil {
		return models.Credentials{}, err
	}

	credentials := jwtToken.mapperUtil.MapUserToCredentials(userToken)
	return credentials, nil

}