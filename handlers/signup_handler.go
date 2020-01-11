package handlers


import (
	"net/http"
	"encoding/json"
	"github.com/cts3njitedu/healthful-heart/services"
	"github.com/cts3njitedu/healthful-heart/models"
)

type SignupHandler struct {
	authenticationService services.IAuthenticationService
	signupService services.ISignupService
}

type ISignupHandler interface {
	GetSignUpPage(w http.ResponseWriter, r *http.Request)
	PostSignUpPage(next http.Handler) http.Handler
}

func NewSingupHandler(authenticationService services.IAuthenticationService, 
	signupService services.ISignupService) *SignupHandler {
		return &SignupHandler{authenticationService, signupService}
}
func (handler *SignupHandler) GetSignUpPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type","application/json")
	page:=handler.authenticationService.GetAuthenticationPage("SIGNUP");
	json.NewEncoder(w).Encode(page)
}


func (handler *SignupHandler) PostSignUpPage(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var page models.Page;
		err := json.NewDecoder(r.Body).Decode(&page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}	
		
	})
}