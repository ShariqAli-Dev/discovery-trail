package types

type TemplateData struct {
	Nonce           string
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}
