package model

type PingRequest struct {
	Data string `json:"data" validate:"required"`
}

type NewAccount struct {
	UserId string `json:"user_id" validate:"required"`
	//ActiveServices   Svc    `json:"active_services" validate:"required"`
	//InactiveServices Svc    `json:"inactive_services" validate:"required"`
}

type UpdateServices struct {
	AccountNumber int    `json:"account_number" validate:"required"`
	ServiceId     string `json:"service_id" validate:"required"`
	UpdateType    string `json:"update_type" validate:"required"`
}
