package model

type PingRequest struct {
	Data string `json:"data" validate:"required"`
}

type NewAccount struct {
	UserId string `json:"user_id" validate:"required"`
}

type UpdateServices struct {
	AccountNumber int    `json:"account_number" validate:"required"`
	ServiceId     string `json:"service_id" validate:"required"`
	UpdateType    string `json:"update_type" validate:"required,oneof=add remove"`
}
type UpdateTransaction struct {
	AccountNumber   int    `json:"account_number" validate:"required"`
	Amount          string `json:"amount" validate:"required"`
	TransactionType string `json:"transaction_type" validate:"required,oneof=debit credit"`
}
