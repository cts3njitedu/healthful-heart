package enrichers

import (
	"github.com/cts3njitedu/healthful-heart/security"
	"github.com/cts3njitedu/healthful-heart/models"
	"errors"
)

type CredentialEnricher struct {
	hasher security.IPasswordHasher
}

func NewCredentialEnricher(hasher security.IPasswordHasher) *CredentialEnricher {
	return &CredentialEnricher{hasher}
}

func (credEnrich *CredentialEnricher) EnrichCredentials (cred *models.Credentials) error {
	pwd, err := credEnrich.hasher.HashPassword(cred.Password);

	if err != nil {
		return errors.New("Enriching credentials went wrong...")
	}

	cred.Password = pwd;
	cred.ConfirmPassword = pwd;

	return nil
}