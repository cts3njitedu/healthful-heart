package handlers

import (
	"net/http"
	"encoding/json"
	"github.com/cts3njitedu/healthful-heart/models"
	"github.com/cts3njitedu/healthful-heart/services"
	"github.com/cts3njitedu/healthful-heart/validators"
	"github.com/cts3njitedu/healthful-heart/security"
	"github.com/cts3njitedu/healthful-heart/repositories/mysqlrepo"
	"github.com/gorilla/context"
	"fmt"
)

type LoginResponse struct {
	IsError bool `json:"isError"`
	Message string `json:"message"`
	Page models.Page `json:"page"`
}

type LoginHandler struct {
	authenticationService services.IAuthenticationService
	loginService services.ILoginService
}

type ILoginHandler interface {
	GetLoginPage(w http.ResponseWriter, r *http.Request)
	PostLoginPage(w http.ResponseWriter, r *http.Request)
}
func NewLoginHandler(authenticationService services.IAuthenticationService, loginService services.ILoginService) *LoginHandler {
	return &LoginHandler{authenticationService, loginService}
}

func (handler *LoginHandler) GetLoginPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type","application/json")
	page:=handler.authenticationService.GetAuthenticationPage("LOGIN");
	json.NewEncoder(w).Encode(page)
}


func (handler *LoginHandler) PostLoginPage(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var page models.Page;
		err := json.NewDecoder(r.Body).Decode(&page)
	
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		newPage, credentials, err := handler.loginService.LoginService(page)

		if err != nil {

			
			
			switch err.(type) {
			case *validators.ValidationError:
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(newPage)
			case *security.PasswordError:
				w.WriteHeader(http.StatusUnauthorized)
				newPage.Errors = append(newPage.Errors, "Invalid username or password")
				json.NewEncoder(w).Encode(LoginResponse{
					IsError: true,
					Page: newPage,
				})
			case *mysqlrepo.UserError: 
				w.WriteHeader(http.StatusUnauthorized)
				newPage.Errors = append(newPage.Errors, "Invalid username or password")
				json.NewEncoder(w).Encode(LoginResponse{
					IsError: true,
					Page: newPage,
				})
			default:
				fmt.Println(err.Error())
				http.Error(w, "Server error code 1", http.StatusInternalServerError)
			}
			
			return
		}
		context.Set(r,"credentials", credentials)
		next.ServeHTTP(w, r)

	})
	

}