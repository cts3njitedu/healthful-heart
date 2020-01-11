package services

import (
	"github.com/cts3njitedu/healthful-heart/models"
)

type IRestructurePageService interface {
	RestructureLoginPage(page *models.Page)
}

type IAuthenticationService interface {
	GetAuthenticationPage(authenticationType string) models.Page
}

type ISignupService interface {
	SignupService(page models.Page) (models.Page, models.Credentials, error)
}