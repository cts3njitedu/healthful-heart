package enrichers

import (
	"github.com/cts3njitedu/healthful-heart/models"
)
type ICredentialEnricher interface {
	EnrichCredentials (cred *models.Credentials) error
}