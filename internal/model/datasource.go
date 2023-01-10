package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

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
	ActiveServices   *Svc      `json:"active_services" sql:"active_services"`
	InactiveServices *Svc      `json:"inactive_services" sql:"inactive_services"`
}
type ColumnUpdate struct {
	UpdateSet string
}
type Svc map[string]struct{}

func (s *Svc) Value() (driver.Value, error) {
	if s == nil || len(*s) == 0 {
		return "{}", nil
	}
	return json.Marshal(s)
}
func (s *Svc) Scan(value any) error {
	if value == nil {
		return nil
	}
	by, ok := value.([]byte)
	if !ok {
		return errors.New("invalid type")
	}
	return json.Unmarshal(by, s)
}

const Schema = `
	(
	user_id varchar(225) not null unique,
	account_number int AUTO_INCREMENT,
	income dec(18,2) DEFAULT 0.00,
	spends dec(18,2) DEFAULT 0.00,
	created_on timestamp not null DEFAULT CURRENT_TIMESTAMP,
	updated_on timestamp not null DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	active_services json,
	inactive_services json,
	primary key (account_number),
	index(user_id)
);
	`
