package enrichers

import (
	"github.com/cts3njitedu/healthful-heart/models"
)

type EnrichExecutor struct {
	enrichers []IEnricher
}

func NewEnrichExecutor(enrichers []IEnricher) *EnrichExecutor {
	return &EnrichExecutor{enrichers}
}
func (execute *EnrichExecutor) Enricher (page *models.Page, pageTemplate models.Page) {
	for _, enricher := range execute.enrichers {
		enricher.Enrich(page, pageTemplate)
	}
}