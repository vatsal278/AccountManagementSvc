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
	Id               string
	AccountNumber    int
	Income           float64
	Spends           float64
	CreatedOn        time.Time
	UpdatedOn        time.Time
	ActiveServices   *Svc
	InactiveServices *Svc
}
type ColumnUpdate struct {
	UpdateSet string
}
type Svc map[string]struct{}

type ColumnUpdate struct {
	UpdateSet string
}

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
