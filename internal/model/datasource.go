package model

import "time"

type DsResponse struct {
	Data string
}

type PingDs struct {
	Data string
}
type Account struct {
	Id               string    `json:"id" validate:"required" sql:"id"`
	AccountNumber    int       `json:"account_number" validate:"required" sql:"account_number"`
	Income           float64   `json:"income" validate:"required" sql:"income"`
	Spends           float64   `json:"spends" validate:"required" sql:"spends"`
	CreatedOn        time.Time `json:"created_on" sql:"created_on"`
	UpdatedOn        time.Time `json:"updated_on" sql:"updated_on"`
	ActiveServices   []string  `json:"active_services" sql:"active_services"`
	InactiveServices []string  `json:"inactive_services" sql:"inactive_services"`
}

const Schema = `
	(
	user_id varchar(225) not null unique,
	account_number int not null unique,
	income int,
	spends int,
	created_on timestamp not null DEFAULT CURRENT_TIMESTAMP,
	updated_on timestamp not null DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	primary key (account_number)
);
	`
