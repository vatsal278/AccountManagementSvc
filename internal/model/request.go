package model

type PingRequest struct {
	Data string `json:"data" validate:"required"`
}

type NewAccount struct {
	UserId string `json:"user_id" validate:"required"`
	//ActiveServices   Svc    `json:"active_services" validate:"required"`
	//InactiveServices Svc    `json:"inactive_services" validate:"required"`
}
