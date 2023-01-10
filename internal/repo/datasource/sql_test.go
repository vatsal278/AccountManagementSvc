package datasource

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	svcCfg "github.com/vatsal278/AccountManagmentSvc/internal/config"
	"github.com/vatsal278/AccountManagmentSvc/internal/model"
	"log"
	"reflect"
	"strings"
	"testing"
	"time"
)

func createTestTable(t *testing.T, db *sql.DB, tableName string, tableStruct string) {
	q := fmt.Sprintf("CREATE Table IF NOT EXISTS %s %s;", tableName, tableStruct)
	_, err := db.Exec(q)
	log.Print(q)
	if err != nil {
		t.Log(err)
	}
}

func deleteTestTable(t *testing.T, db *sql.DB, tableName string) {
	q := fmt.Sprintf("DROP TABLE %s", tableName)
	_, err := db.Exec(q)
	if err != nil {
		t.Fatal(err.Error())
	}
}
func TestSqlDs_HealthCheck(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing due to unavailability of testing environment")
	}
	//include a failure case
	dbcfg := svcCfg.DbCfg{
		Port:      "9085",
		Host:      "localhost",
		Driver:    "mysql",
		User:      "root",
		Pass:      "pass",
		DbName:    "useracc",
		TableName: "newTemp",
	}
	dataBase := svcCfg.Connect(dbcfg, dbcfg.TableName)
	svcConfig := svcCfg.SvcConfig{
		DbSvc: svcCfg.DbSvc{Db: dataBase},
	}
	dB := NewSql(svcCfg.DbSvc(svcConfig.DbSvc), "newTemp")

	tests := []struct {
		name        string
		setupFunc   func()
		cleanupFunc func()
		filter      map[string]interface{}
		validator   func(bool)
		dbInterface DataSourceI
	}{
		{
			name: "SUCCESS::Health check",
			validator: func(res bool) {
				if res != true {
					t.Errorf("Want: %v, Got: %v", true, res)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res := dB.HealthCheck()

			if tt.validator != nil {
				tt.validator(res)
			}
		})
	}
}
func TestGet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing due to unavailability of testing environment")
	}
	dbcfg := svcCfg.DbCfg{
		Port:      "9085",
		Host:      "localhost",
		Driver:    "mysql",
		User:      "root",
		Pass:      "pass",
		DbName:    "accmgmt",
		TableName: "newTemp",
	}
	dataBase := svcCfg.Connect(dbcfg, dbcfg.TableName)
	svcConfig := svcCfg.SvcConfig{
		DbSvc: svcCfg.DbSvc{Db: dataBase},
		Cfg:   &svcCfg.Config{DataBase: dbcfg},
	}
	dB := sqlDs{
		sqlSvc: svcConfig.DbSvc.Db,
		table:  svcConfig.Cfg.DataBase.TableName,
	}

	tests := []struct {
		name        string
		setupFunc   func()
		cleanupFunc func()
		filter      map[string]interface{}
		validator   func([]model.Account, error)
		dbInterface DataSourceI
	}{
		{
			name: "SUCCESS::Get",
			filter: map[string]interface{}{
				"user_id": "1234",
			},
			setupFunc: func() {
				dataBase.Exec("DROP TABLE newTemp")
				createTestTable(t, dataBase, "newTemp", model.Schema)
				err := dB.Insert(model.Account{
					Id:             "1234",
					ActiveServices: nil,
					InactiveServices: &model.Svc{
						"1": {},
					},
				})
				if err != nil {
					t.Fatal(err.Error())
				}
			},
			cleanupFunc: func() {
				deleteTestTable(t, dataBase, "newTemp")
			},
			validator: func(rows []model.Account, err error) {
				temp := model.Account{
					Id:            "1234",
					AccountNumber: 1,
					CreatedOn:     time.Now(),
				}

				if !reflect.DeepEqual(rows[0].Id, temp.Id) {
					t.Errorf("Want: %v, Got: %v", temp.Id, rows[0].Id)
				}
				if !reflect.DeepEqual(rows[0].AccountNumber, temp.AccountNumber) {
					t.Errorf("Want: %v, Got: %v", temp.AccountNumber, rows[0].AccountNumber)
				}
				if !reflect.DeepEqual(rows[0].Income, temp.Income) {
					t.Errorf("Want: %v, Got: %v", temp.Income, rows[0].Income)
				}
				if !reflect.DeepEqual(rows[0].Spends, temp.Spends) {
					t.Errorf("Want: %v, Got: %v", temp.Spends, rows[0].Spends)
				}
				if !reflect.DeepEqual(rows[0].Income, temp.Income) {
					t.Errorf("Want: %v, Got: %v", temp.Income, rows[0].Income)
				}
				if !reflect.DeepEqual(rows[0].ActiveServices, &model.Svc{}) {
					t.Errorf("Want: %v, Got: %v", temp.Income, rows[0].Income)
				}
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err)
				}
			},
		},
		{
			name: "SUCCESS::Get:: multiple articles",
			filter: map[string]interface{}{
				"user_id": "1234",
			},
			setupFunc: func() {
				dataBase.Exec("DROP TABLE newTemp")
				createTestTable(t, dataBase, "newTemp", model.Schema)
				err := dB.Insert(model.Account{
					Id:             "1234",
					ActiveServices: nil,
					InactiveServices: &model.Svc{
						"1": {},
					},
				})
				if err != nil {
					t.Fatal(err.Error())
				}
				err = dB.Insert(model.Account{
					Id: "4321",
					ActiveServices: &model.Svc{
						"2": {},
					},
					InactiveServices: &model.Svc{
						"1": {},
					},
				})
				if err != nil {
					t.Fatal(err.Error())
				}
			},
			cleanupFunc: func() {
				deleteTestTable(t, dataBase, "newTemp")
			},
			validator: func(rows []model.Account, err error) {
				temp := model.Account{
					Id:               "1234",
					AccountNumber:    1,
					InactiveServices: &model.Svc{"1": {}},
				}
				if !reflect.DeepEqual(rows[0].Id, temp.Id) {
					t.Errorf("Want: %v, Got: %v", temp.Id, rows[0].Id)
				}
				if !reflect.DeepEqual(rows[0].AccountNumber, temp.AccountNumber) {
					t.Errorf("Want: %v, Got: %v", temp.AccountNumber, rows[0].AccountNumber)
				}
				if !reflect.DeepEqual(rows[0].Income, temp.Income) {
					t.Errorf("Want: %v, Got: %v", temp.Income, rows[0].Income)
				}
				if !reflect.DeepEqual(rows[0].Spends, temp.Spends) {
					t.Errorf("Want: %v, Got: %v", temp.Spends, rows[0].Spends)
				}
				if !reflect.DeepEqual(rows[0].InactiveServices, temp.InactiveServices) {
					t.Errorf("Want: %v, Got: %v", temp.InactiveServices, rows[0].InactiveServices)
				}
				if !reflect.DeepEqual(rows[0].ActiveServices, &model.Svc{}) {
					t.Errorf("Want: %v, Got: %v", temp.Income, rows[0].Income)
				}
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err)
				}
			},
		},
		{
			name: "SUCCESS::Get::no user found",
			filter: map[string]interface{}{
				"account_number": 1,
			},
			setupFunc: func() {
				createTestTable(t, dataBase, "newTemp", model.Schema)
			},
			cleanupFunc: func() {
				deleteTestTable(t, dataBase, "newTemp")
			},
			validator: func(rows []model.Account, err error) {
				if len(rows) != 0 {
					t.Errorf("Want: %v, Got: %v", 0, len(rows))
				}
			},
		},
		{
			name: "failure::Get::scan error", //scan should return an error
			filter: map[string]interface{}{
				"user_id": "12345",
			},
			setupFunc: func() {
				dataBase.Exec("DROP TABLE newTemp")
				createTestTable(t, dataBase, "newTemp", `
	(
	user_id varchar(225) not null unique,
	account_number int AUTO_INCREMENT,
	income text,
	spends dec(18,2) DEFAULT 0.00,
	created_on int DEFAULT 0,
	updated_on timestamp not null DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	active_services json,
	inactive_services json,
	primary key (account_number),
	index(user_id)
);
	`)

				_, err := dataBase.Exec("INSERT INTO newTemp(user_id, created_on) VALUES(?,?)", "12345", 1)
				if err != nil {
					t.Error(err.Error())
					return
				}
				if err != nil {
					t.Fatal(err.Error())
				}

			},
			cleanupFunc: func() {
				deleteTestTable(t, dataBase, "newTemp")
			},
			validator: func(rows []model.Account, err error) {
				if !strings.Contains(err.Error(), "sql: Scan error on column") {
					t.Errorf("Want: %v, Got: %v", "sql: Scan error on column", err.Error())
				}
			},
		},
		{
			name:   "FAILURE:: query error",
			filter: map[string]interface{}{"userid": "v@mail.com"},
			setupFunc: func() {
				//dataBase.Exec("DROP TABLE newTemp")
				createTestTable(t, dataBase, "newTemp", model.Schema)

			},
			cleanupFunc: func() {
				tableName := "newTemp"
				deleteTestTable(t, dataBase, tableName)

			},
			validator: func(rows []model.Account, err error) {
				if !strings.Contains(err.Error(), "Unknown column") {
					t.Errorf("Want: %v, Got: %v", "Unknown column", err)
				}
			},
		},
	}

	// to execute the tests in the table
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// STEP 1: seting up all instances for the specific test case
			if tt.setupFunc != nil {
				tt.setupFunc()
			}

			// STEP 2: call the test function
			rows, err := dB.Get(tt.filter)
			t.Log(rows)

			// STEP 3: validation of output
			if tt.validator != nil {
				tt.validator(rows, err)
			}

			// STEP 4: clean up/remove up all instances for the specific test case
			if tt.cleanupFunc != nil {
				tt.cleanupFunc()
			}
		})
	}
}

//
func TestInsert(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing due to unavailability of testing environment")
	}
	//add short flag for dbtest case
	dbcfg := svcCfg.DbCfg{
		Port:      "9085",
		Host:      "localhost",
		Driver:    "mysql",
		User:      "root",
		Pass:      "pass",
		DbName:    "accmgmt",
		TableName: "newTemp",
	}
	dataBase := svcCfg.Connect(dbcfg, dbcfg.TableName)
	svcConfig := svcCfg.SvcConfig{
		Cfg:   &svcCfg.Config{DataBase: dbcfg},
		DbSvc: svcCfg.DbSvc{Db: dataBase},
	}
	dB := sqlDs{
		sqlSvc: svcConfig.DbSvc.Db,
		table:  svcConfig.Cfg.DataBase.TableName,
	}
	// table driven tests
	tests := []struct {
		name        string
		tableName   string
		data        model.Account
		setupFunc   func()
		cleanupFunc func()
		filter      map[string]interface{}
		validator   func(error)
	}{
		{
			name: "SUCCESS:: Insert Article",
			data: model.Account{
				Id:               "12345",
				ActiveServices:   &model.Svc{"1": {}},
				InactiveServices: &model.Svc{"2": {}},
			},
			setupFunc: func() {
				tableName := "newTemp"
				createTestTable(t, dataBase, tableName, model.Schema)
			},
			cleanupFunc: func() {
				tableName := "newTemp"
				deleteTestTable(t, dataBase, tableName)
			},
			validator: func(err error) {
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
				user, err := dB.Get(map[string]interface{}{"user_id": "12345"})
				if err != nil {
					t.Errorf("unable to get data from db")
				}
				temp := model.Account{
					Id:               "12345",
					AccountNumber:    1,
					ActiveServices:   &model.Svc{"1": {}},
					InactiveServices: &model.Svc{"2": {}},
				}
				if !reflect.DeepEqual(user[0].Id, temp.Id) {
					t.Errorf("Want: %v, Got: %v", temp.Id, user[0].Id)
				}
				if !reflect.DeepEqual(user[0].ActiveServices, temp.ActiveServices) {
					t.Errorf("Want: %v, Got: %v", temp.ActiveServices, user[0].ActiveServices)
				}
				if !reflect.DeepEqual(user[0].InactiveServices, temp.InactiveServices) {
					t.Errorf("Want: %v, Got: %v", temp.InactiveServices, user[0].InactiveServices)
				}
				if !reflect.DeepEqual(user[0].AccountNumber, temp.AccountNumber) {
					t.Errorf("Want: %v, Got: %v", temp.AccountNumber, user[0].AccountNumber)
				}
				if !reflect.DeepEqual(user[0].Income, temp.Income) {
					t.Errorf("Want: %v, Got: %v", temp.Income, user[0].Income)
				}
				if !reflect.DeepEqual(user[0].Spends, temp.Spends) {
					t.Errorf("Want: %v, Got: %v", temp.Spends, user[0].Spends)
				}
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err)
				}
			},
		},
		{
			name: "SUCCESS:: Insert Article:: Insert Article when data already present",
			data: model.Account{
				Id:               "12345",
				ActiveServices:   &model.Svc{"1": {}},
				InactiveServices: &model.Svc{"2": {}},
			},
			setupFunc: func() {
				tableName := "newTemp"
				createTestTable(t, dataBase, tableName, model.Schema)
				err := dB.Insert(model.Account{
					Id:               "01",
					ActiveServices:   nil,
					InactiveServices: nil,
				})
				if err != nil {
					t.Fatal(err)
				}
			},
			cleanupFunc: func() {
				tableName := "newTemp"
				deleteTestTable(t, dataBase, tableName)
			},
			validator: func(err error) {
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
				user, err := dB.Get(map[string]interface{}{"user_id": "01"})
				if err != nil {
					t.Errorf("unable to get data from db")
				}
				temp := model.Account{
					Id:               "01",
					AccountNumber:    1,
					ActiveServices:   &model.Svc{},
					InactiveServices: &model.Svc{},
				}
				if !reflect.DeepEqual(user[0].Id, temp.Id) {
					t.Errorf("Want: %v, Got: %v", temp.Id, user[0].Id)
				}
				if !reflect.DeepEqual(user[0].ActiveServices, temp.ActiveServices) {
					t.Errorf("Want: %v, Got: %v", temp.ActiveServices, user[0].ActiveServices)
				}
				if !reflect.DeepEqual(user[0].InactiveServices, temp.InactiveServices) {
					t.Errorf("Want: %v, Got: %v", temp.InactiveServices, user[0].InactiveServices)
				}
				if !reflect.DeepEqual(user[0].AccountNumber, temp.AccountNumber) {
					t.Errorf("Want: %v, Got: %v", temp.AccountNumber, user[0].AccountNumber)
				}
				if !reflect.DeepEqual(user[0].Income, temp.Income) {
					t.Errorf("Want: %v, Got: %v", temp.Income, user[0].Income)
				}
				if !reflect.DeepEqual(user[0].Spends, temp.Spends) {
					t.Errorf("Want: %v, Got: %v", temp.Spends, user[0].Spends)
				}
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err)
				}
			},
		},
		{
			name: "FAILURE:: column mismatch",
			data: model.Account{
				ActiveServices:   &model.Svc{"1": {}},
				InactiveServices: &model.Svc{"2": {}},
			},
			setupFunc: func() {
				dataBase.Exec("DROP TABLE newTemp")
				tableName := "newTemp"
				createTestTable(t, dataBase, tableName, `
	(
	user_id int not null unique,
	account_number int AUTO_INCREMENT,
	income dec(18,2) DEFAULT 0.00,
	spends dec(18,2) DEFAULT 0.00,
	created_on timestamp not null DEFAULT CURRENT_TIMESTAMP,
	updated_on timestamp not null DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	active_services int,
	inactive_services int,
	primary key (account_number),
	index(user_id)
);
	`)
			},
			cleanupFunc: func() {
				tableName := "newTemp"
				deleteTestTable(t, dataBase, tableName)
			},
			validator: func(err error) {
				var tempErr = errors.New("Error 1366 (HY000): Incorrect integer value: for column 'active_services' at row 1")
				if !strings.Contains(err.Error(), "Error 1366") {
					t.Errorf("Want: %v, Got: %v", tempErr, err)
				}
			},
		},
	}
	// to execute the tests in the table
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFunc != nil {
				tt.setupFunc()
			}
			// STEP 2: call the test function
			err := dB.Insert(tt.data)
			// STEP 3: validation of output
			if tt.validator != nil {
				tt.validator(err)
			}
			// STEP 4: clean up/remove up all instances for the specific test case
			if tt.cleanupFunc != nil {
				tt.cleanupFunc()
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing due to unavailability of testing environment")
	}
	dbcfg := svcCfg.DbCfg{
		Port:      "9085",
		Host:      "localhost",
		Driver:    "mysql",
		User:      "root",
		Pass:      "pass",
		DbName:    "accmgmt",
		TableName: "newTemp",
	}
	dataBase := svcCfg.Connect(dbcfg, dbcfg.TableName)
	svcConfig := svcCfg.SvcConfig{
		Cfg:   &svcCfg.Config{DataBase: dbcfg},
		DbSvc: svcCfg.DbSvc{Db: dataBase},
	}
	dB := sqlDs{
		sqlSvc: svcConfig.DbSvc.Db,
		table:  svcConfig.Cfg.DataBase.TableName,
	}
	// table driven tests
	tests := []struct {
		name        string
		tableName   string
		dataSet     map[string]interface{}
		dataWhere   map[string]interface{}
		setupFunc   func()
		cleanupFunc func()
		filter      map[string]interface{}
		validator   func(error)
	}{
		{
			name:      "SUCCESS:: Update",
			dataSet:   map[string]interface{}{"active_services": model.Svc{"1": {}, "2": {}}},
			dataWhere: map[string]interface{}{"user_id": "1234"},
			setupFunc: func() {
				tableName := "accdatabase"
				createTestTable(t, dataBase, tableName, model.Schema)
				err := dB.Insert(model.Account{
					Id: "123",
					//ActiveServices:   &model.Svc{"1": {}},
					//InactiveServices: &model.Svc{"2": {}},
				})
				if err != nil {
					t.Fatal(err)
				}
			},
			cleanupFunc: func() {
				tableName := "newTemp"
				deleteTestTable(t, dataBase, tableName)
			},
			validator: func(err error) {
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
				user, err := dB.Get(map[string]interface{}{"user_id": "1234"})
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}

			},
		},
		{
			name:      "Failure:: Update",
			dataSet:   map[string]interface{}{"active_services": model.Svc{"3": {}}},
			dataWhere: map[string]interface{}{"user_id": "1234"},
			setupFunc: func() {
				tableName := "newTemp"
				createTestTable(t, dataBase, tableName, model.Schema)
				err := dB.Insert(model.Account{
					Id:               "1234",
					ActiveServices:   &model.Svc{"1": {}},
					InactiveServices: &model.Svc{"2": {}},
				})
				if err != nil {
					t.Fatal(err)
				}
			},
			cleanupFunc: func() {
				tableName := "newTemp"
				deleteTestTable(t, dataBase, tableName)
			},
			validator: func(err error) {
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
				user, err := dB.Get(map[string]interface{}{"user_id": "1234"})
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
				t.Log(user[0].ActiveServices)
			},
		},
	}
	// to execute the tests in the table
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFunc != nil {
				tt.setupFunc()
			}
			// STEP 2: call the test function
			err := dB.Update(tt.dataSet, tt.dataWhere)

			// STEP 3: validation of output
			if tt.validator != nil {
				tt.validator(err)
			}

			// STEP 4: clean up/remove up all instances for the specific test case
			if tt.cleanupFunc != nil {
				tt.cleanupFunc()
			}
		})
	}
}
