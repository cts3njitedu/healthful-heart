package enrichers

import (
	"github.com/cts3njitedu/healthful-heart/models"
)
type ICredentialEnricher interface {
	EnrichCredentials (cred *models.Credentials) error
}

type IEnricher interface {
	Enrich(page *models.Page, pageTemplate models.Page)
}

type IEnricherExecutor interface {
	Enricher(page *models.Page, pageTemplate models.Page)
}