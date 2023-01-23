package model

type AccountSummary struct {
	AccountNumber    int     `json:"account_number,omitempty"`
	Income           float64 `json:"income"`
	Spends           float64 `json:"spends"`
	ActiveServices   *Svc    `json:"active_services"`
	InactiveServices *Svc    `json:"inactive_services"`
}
type CacheResponse struct {
	Status      int
	Response    string
	ContentType string
}
