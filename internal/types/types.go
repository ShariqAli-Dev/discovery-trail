package types

import "github.com/shariqali-dev/discovery-trail/internal/models"

type TemplateData struct {
	Nonce           string
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
	Account         models.Account
}
