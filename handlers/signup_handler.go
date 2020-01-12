package handlers


import (
	"net/http"
	"encoding/json"
	"github.com/cts3njitedu/healthful-heart/services"
	"github.com/cts3njitedu/healthful-heart/models"
	"github.com/cts3njitedu/healthful-heart/validators"
	"github.com/gorilla/context"
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
		newPage, credentials, err := handler.signupService.SignupService(page)

		if err != nil {

			
			switch err.(type) {
			case *validators.ValidationError:
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(newPage)
			default:
				http.Error(w, "Server error code 1", http.StatusInternalServerError)
				
			}
			
			return
		}
		context.Set(r,"credentials", credentials)
		next.ServeHTTP(w, r)
	})
}