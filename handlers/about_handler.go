package handlers

import (
	"net/http"
	"fmt"
	"github.com/gorilla/context"
)

type AboutHandler struct {

}

type IAboutHandler interface {
	GetAboutPage(w http.ResponseWriter, r *http.Request)
}

func NewAboutHandler() *AboutHandler {
	return &AboutHandler{}
}


func (handler *AboutHandler) GetAboutPage(w http.ResponseWriter, r *http.Request) {
	claims:=context.Get(r,"credentials")
	fmt.Printf("Decoded claims %+v\n", claims)
	fmt.Println("In about page")
	w.Write([]byte(`{"message": "In about page"}`))

}