package datasource

import (
	"database/sql"
	"fmt"
	"github.com/vatsal278/AccountManagmentSvc/internal/config"
	"github.com/vatsal278/AccountManagmentSvc/internal/model"
	"strings"
)

type sqlDs struct {
	sqlSvc *sql.DB
	table  string
}

func NewSql(dbSvc config.DbSvc, tableName string) DataSourceI {
	return &sqlDs{
		sqlSvc: dbSvc.Db,
		table:  tableName,
	}
}

func (d sqlDs) HealthCheck() bool {
	err := d.sqlSvc.Ping()
	if err != nil {
		return false
	}
	return true
}

func (d sqlDs) Get(filter map[string]interface{}) ([]model.Account, error) {
	//order the queries based on email address
	var user model.Account
	var users []model.Account
	q := fmt.Sprintf("SELECT user_id, account_number, income, spends, registered_on, updated_on, active_services, inactive_services FROM %s", d.table)

	filterClause := []string{}

	for k, v := range filter {
		switch v.(type) {
		case string:
			filterClause = append(filterClause, fmt.Sprintf("%s = '%s'", k, v))
		default:
			filterClause = append(filterClause, fmt.Sprintf("%s = %+v", k, v))
		}
	}
	if len(filterClause) > 0 {
		q += fmt.Sprintf(" WHERE %s", strings.Join(filterClause, " AND "))
	}

	q += " ORDER BY email;"
	rows, err := d.sqlSvc.Query(q)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&user.Id, &user.AccountNumber, &user.Income, &user.Spends, &user.CreatedOn, &user.UpdatedOn, &user.ActiveServices, &user.InactiveServices)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (d sqlDs) Insert(user model.Account) error {
	queryString := fmt.Sprintf("INSERT INTO %s", d.table)
	_, err := d.sqlSvc.Exec(queryString+"(user_id, registered_on, active_services, inactive_services) VALUES(?, ?,?,?)", user.Id, user.CreatedOn, user.ActiveServices, user.InactiveServices)
	if err != nil {
		return err
	}
	return err
}

func (d sqlDs) Update(filterSet map[string]interface{}, filterWhere map[string]interface{}) error {
	queryString := fmt.Sprintf("UPDATE %s ", d.table)
	filterClause := []string{}

	for k, v := range filterSet {
		switch v.(type) {
		case string:
			filterClause = append(filterClause, fmt.Sprintf("%s = '%+v'", k, v))
		default:
			filterClause = append(filterClause, fmt.Sprintf("%s = %+v", k, v))
		}
	}
	if len(filterClause) > 0 {
		queryString += fmt.Sprintf(" SET %s", strings.Join(filterClause, " , "))
	}
	filterClauseWhere := []string{}

	for k, v := range filterWhere {
		switch v.(type) {
		case string:
			filterClauseWhere = append(filterClauseWhere, fmt.Sprintf("%s = '%+v'", k, v))
		default:
			filterClauseWhere = append(filterClauseWhere, fmt.Sprintf("%s = %+v", k, v))
		}
	}
	if len(filterClauseWhere) > 0 {
		queryString += fmt.Sprintf(" WHERE %s", strings.Join(filterClauseWhere, " AND "))
	}

	queryString += " ;"
	_, err := d.sqlSvc.Exec(queryString)
	if err != nil {
		return err
	}
	return nil

}
