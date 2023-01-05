package datasource

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	svcCfg "github.com/vatsal278/AccountManagmentSvc/internal/config"
	"github.com/vatsal278/AccountManagmentSvc/internal/model"
	"reflect"
	"testing"
	"time"
)

func createTestTable(t *testing.T, db *sql.DB, tableName string, tableStruct string) {
	q := fmt.Sprintf("CREATE Table IF NOT EXISTS %s (%s);", tableName, tableStruct)
	_, err := db.Exec(q)
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
		Port:      "9095",
		Host:      "localhost",
		Driver:    "mysql",
		User:      "root",
		Pass:      "pass",
		DbName:    "usermgmt",
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
				"account_number": "1",
			},
			setupFunc: func() {
				dataBase.Exec("DROP TABLE newTemp")
				createTestTable(t, dataBase, "newTemp", model.Schema)
				err := dB.Insert(model.Account{
					Id:        "1234",
					CreatedOn: time.Now(),
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
					Id:        "1234",
					CreatedOn: time.Now(),
				}
				if !reflect.DeepEqual(rows[0], temp) {
					t.Errorf("Want: %v, Got: %v", temp, rows[0])
				}
				//if !reflect.DeepEqual(rows[0].Name, temp.Name) {
				//	t.Errorf("Want: %v, Got: %v", temp.Name, rows[0].Name)
				//}
				//if !reflect.DeepEqual(rows[0].Company, temp.Company) {
				//	t.Errorf("Want: %v, Got: %v", temp.Company, rows[0].Company)
				//}
				//if err != nil {
				//	t.Errorf("Want: %v, Got: %v", nil, err)
				//}
			},
		},
		//{
		//	name: "SUCCESS::Get:: multiple articles",
		//	filter: map[string]interface{}{
		//		"name":   "vatsal",
		//		"active": false,
		//	},
		//	setupFunc: func() {
		//		createTestTable(t, dataBase, "newTemp", "user_id varchar(225) not null, email varchar(225) not null unique, company_name varchar(225), name varchar(225) not null, password varchar(225) not null DEFAULT 00000000, registered_on timestamp not null DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, updated_on timestamp not null DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, active boolean not null default false, active_devices int(50) not null default 0, salt varchar(225) not null default 0000, primary key (email) ")
		//		err := dB.Insert(model.User{
		//			Email:    "v@mail.com",
		//			Company:  "company",
		//			Password: "pass",
		//			Name:     "vatsal",
		//
		//			RegisteredOn: time.Now(),
		//		})
		//		if err != nil {
		//			t.Fatal(err.Error())
		//		}
		//		err = dB.Insert(model.User{
		//			Email:    "a@mail.com",
		//			Company:  "company",
		//			Password: "pass",
		//			Name:     "vatsal",
		//
		//			RegisteredOn: time.Now(),
		//		})
		//		if err != nil {
		//			t.Fatal(err.Error())
		//		}
		//	},
		//	cleanupFunc: func() {
		//		deleteTestTable(t, dataBase, "newTemp")
		//	},
		//	validator: func(rows []model.User, err error) {
		//		temp := model.User{
		//			Email:    "v@mail.com",
		//			Password: "pass",
		//			Company:  "company",
		//			Name:     "vatsal",
		//		}
		//		if !reflect.DeepEqual(rows[1].Email, temp.Email) {
		//			t.Errorf("Want: %v, Got: %v", temp.Email, rows[1].Email)
		//		}
		//		if !reflect.DeepEqual(rows[1].Name, temp.Name) {
		//			t.Errorf("Want: %v, Got: %v", temp.Name, rows[1].Name)
		//		}
		//		if !reflect.DeepEqual(rows[1].Company, temp.Company) {
		//			t.Errorf("Want: %v, Got: %v", temp.Company, rows[1].Company)
		//		}
		//		if err != nil {
		//			t.Errorf("Want: %v, Got: %v", nil, err)
		//		}
		//	},
		//},
		//{
		//	name: "SUCCESS::Get::no user found",
		//	filter: map[string]interface{}{
		//		"email": "v@mail.com",
		//	},
		//	setupFunc: func() {
		//		createTestTable(t, dataBase, "newTemp", "user_id varchar(225) not null, email varchar(225) not null unique, company_name varchar(225), name int(225) not null, password varchar(225) not null DEFAULT 00000000, registered_on timestamp not null DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, updated_on timestamp not null DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, active boolean not null default false, active_devices int(50) not null default 0, salt varchar(225) not null default 0000, primary key (email) ")
		//	},
		//	cleanupFunc: func() {
		//		deleteTestTable(t, dataBase, "newTemp")
		//	},
		//	validator: func(rows []model.User, err error) {
		//		if len(rows) != 0 {
		//			t.Errorf("Want: %v, Got: %v", 0, len(rows))
		//		}
		//	},
		//},
		//{
		//	name: "failure::Get::scan error", //scan should return an error
		//	filter: map[string]interface{}{
		//		"email": "vmail.com",
		//	},
		//	setupFunc: func() {
		//		createTestTable(t, dataBase, "newTemp", "user_id varchar(225) not null, email varchar(225) not null unique, company_name varchar(225), name int(225) not null, password varchar(225) not null DEFAULT 00000000, registered_on timestamp not null DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, updated_on timestamp not null DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, active boolean not null default false, active_devices int(50) not null default 0, salt varchar(225) not null default 0000, primary key (email) ")
		//		_, err := dataBase.Exec("INSERT INTO newTemp(user_id, name, registered_on, email, salt, company_name, password,active ) VALUES(?,?, ?,?,?,?, ?,?)", "123", "1", time.Now(), "vmail.com", 1, nil, "wxyz", false)
		//		if err != nil {
		//			t.Error(err.Error())
		//			return
		//		}
		//		if err != nil {
		//			t.Fatal(err.Error())
		//		}
		//
		//	},
		//	cleanupFunc: func() {
		//		deleteTestTable(t, dataBase, "newTemp")
		//	},
		//	validator: func(rows []model.User, err error) {
		//		if !strings.Contains(err.Error(), "sql: Scan error on column") {
		//			t.Errorf("Want: %v, Got: %v", "sql: Scan error on column", err.Error())
		//		}
		//	},
		//},
		//{
		//	name:   "FAILURE:: query error",
		//	filter: map[string]interface{}{"email": "v@mail.com"},
		//	setupFunc: func() {
		//		dataBase.Exec("DROP TABLE newTemp")
		//		createTestTable(t, dataBase, "newTemp", "active boolean not null default false")
		//	},
		//	cleanupFunc: func() {
		//		tableName := "newTemp"
		//		deleteTestTable(t, dataBase, tableName)
		//
		//	},
		//	validator: func(rows []model.User, err error) {
		//		if len(rows) != 0 {
		//			t.Errorf("Want: %v, Got: %v", 0, len(rows))
		//		}
		//		var tempErr = errors.New("Error 1054: Unknown column 'user_id' in 'field list'")
		//		if !reflect.DeepEqual(tempErr.Error(), err.Error()) {
		//			t.Errorf("Want: %v, Got: %v", tempErr, err)
		//		}
		//	},
		//},
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
//func TestInsert(t *testing.T) {
//	if testing.Short() {
//		t.Skip("skipping testing due to unavailability of testing environment")
//	}
//	//add short flag for dbtest case
//	dbcfg := config.DbCfg{
//		Port:      "9095",
//		Host:      "localhost",
//		Driver:    "mysql",
//		User:      "root",
//		Pass:      "pass",
//		DbName:    "usermgmt",
//		TableName: "newTemp",
//	}
//	dataBase := config.Connect(dbcfg, dbcfg.TableName)
//	svcConfig := config.SvcConfig{
//		Cfg:   &config.Config{DataBase: dbcfg},
//		DbSvc: config.DbSvc{Db: dataBase},
//	}
//	dB := sqlDs{
//		sqlSvc: svcConfig.DbSvc.Db,
//		table:  svcConfig.Cfg.DataBase.TableName,
//	}
//	// table driven tests
//	tests := []struct {
//		name        string
//		tableName   string
//		data        model.User
//		setupFunc   func()
//		cleanupFunc func()
//		filter      map[string]interface{}
//		validator   func(error)
//	}{
//		{
//			name: "SUCCESS:: Insert Article",
//			data: model.User{
//				Email:    "v@mail.com",
//				Company:  "company",
//				Password: "pass",
//				Name:     "vatsal",
//
//				RegisteredOn: time.Now(),
//			},
//			setupFunc: func() {
//				tableName := "newTemp"
//				createTestTable(t, dataBase, tableName, "user_id varchar(225) not null, email varchar(225) not null unique, company_name varchar(225), name varchar(225) not null, password varchar(225) not null DEFAULT 00000000, registered_on timestamp not null DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, updated_on timestamp not null DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, active boolean not null default false, active_devices int(50) not null default 0, salt varchar(225) not null default 0000, primary key (email) ")
//			},
//			cleanupFunc: func() {
//				tableName := "newTemp"
//				deleteTestTable(t, dataBase, tableName)
//			},
//			validator: func(err error) {
//				if err != nil {
//					t.Errorf("Want: %v, Got: %v", nil, err.Error())
//				}
//				user, err := dB.Get(map[string]interface{}{"email": "v@mail.com"})
//				if err != nil {
//					t.Errorf("unable to get data from db")
//				}
//				temp := model.User{
//					Email:    "v@mail.com",
//					Company:  "company",
//					Password: "pass",
//					Name:     "vatsal"}
//				if !reflect.DeepEqual(user[0].Name, temp.Name) {
//					t.Errorf("Want: %v, Got: %v", temp.Name, user[0].Name)
//				}
//				if !reflect.DeepEqual(user[0].Company, temp.Company) {
//					t.Errorf("Want: %v, Got: %v", temp.Company, user[0].Company)
//				}
//				if err != nil {
//					t.Errorf("Want: %v, Got: %v", nil, err)
//				}
//			},
//		},
//		{
//			name: "SUCCESS:: Insert Article when data already present",
//			data: model.User{
//				Email:    "v@mail.com",
//				Company:  "company",
//				Password: "pass",
//				Name:     "vatsal",
//
//				RegisteredOn: time.Now(),
//			},
//			setupFunc: func() {
//				tableName := "newTemp"
//				createTestTable(t, dataBase, tableName, "user_id varchar(225) not null, email varchar(225) not null unique, company_name varchar(225), name varchar(225) not null, password varchar(225) not null DEFAULT 00000000, registered_on timestamp not null DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, updated_on timestamp not null DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, active boolean not null default false, active_devices int(50) not null default 0, salt varchar(225) not null default 0000, primary key (email) ")
//				err := dB.Insert(model.User{
//					Email:    "vm",
//					Company:  "company1",
//					Password: "pass1",
//					Name:     "vatsal1",
//
//					RegisteredOn: time.Now(),
//				})
//				if err != nil {
//					t.Fatal(err)
//				}
//			},
//			cleanupFunc: func() {
//				tableName := "newTemp"
//				deleteTestTable(t, dataBase, tableName)
//
//			},
//			validator: func(err error) {
//				if err != nil {
//					t.Errorf("Want: %v, Got: %v", nil, err.Error())
//				}
//				user, err := dB.Get(map[string]interface{}{"email": "v@mail.com"})
//				if err != nil {
//					t.Errorf("unable to get data from db")
//				}
//				temp := model.User{
//					Email:    "v@mail.com",
//					Company:  "company",
//					Password: "pass",
//					Name:     "vatsal"}
//				if !reflect.DeepEqual(user[0].Name, temp.Name) {
//					t.Errorf("Want: %v, Got: %v", temp.Name, user[0].Name)
//				}
//				if !reflect.DeepEqual(user[0].Company, temp.Company) {
//					t.Errorf("Want: %v, Got: %v", temp.Company, user[0].Company)
//				}
//				if err != nil {
//					t.Errorf("Want: %v, Got: %v", nil, err)
//				}
//			},
//		},
//		{
//			name: "FAILURE:: column mismatch",
//			data: model.User{
//				Email:    "vm",
//				Company:  "company",
//				Password: "pass",
//				Name:     "vatsal",
//
//				RegisteredOn: time.Now(),
//			},
//			setupFunc: func() {
//				tableName := "newTemp"
//				createTestTable(t, dataBase, tableName, "user_id varchar(225) not null, email varchar(225) not null unique, company_name int(225), name varchar(225) not null, password varchar(225) not null DEFAULT 00000000, registered_on timestamp not null DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, updated_on timestamp not null DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, active boolean not null default false, active_devices int(50) not null default 0, salt varchar(225) not null default 0000, primary key (email) ")
//			},
//			cleanupFunc: func() {
//				tableName := "newTemp"
//				deleteTestTable(t, dataBase, tableName)
//
//			},
//			validator: func(err error) {
//				var tempErr = errors.New("Error 1366: Incorrect integer value: 'company' for column 'company_name' at row 1")
//
//				if !strings.Contains(err.Error(), "Error 1366") {
//					t.Errorf("Want: %v, Got: %v", tempErr, err)
//				}
//			},
//		},
//	}
//	// to execute the tests in the table
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if tt.setupFunc != nil {
//				tt.setupFunc()
//			}
//			// STEP 2: call the test function
//			err := dB.Insert(tt.data)
//			// STEP 3: validation of output
//			if tt.validator != nil {
//				tt.validator(err)
//			}
//			// STEP 4: clean up/remove up all instances for the specific test case
//			if tt.cleanupFunc != nil {
//				tt.cleanupFunc()
//			}
//		})
//	}
//}
//
//func TestUpdate(t *testing.T) {
//	if testing.Short() {
//		t.Skip("skipping testing due to unavailability of testing environment")
//	}
//	dbcfg := config.DbCfg{
//		Port:      "9095",
//		Host:      "localhost",
//		Driver:    "mysql",
//		User:      "root",
//		Pass:      "pass",
//		DbName:    "usermgmt",
//		TableName: "newTemp",
//	}
//	dataBase := config.Connect(dbcfg, dbcfg.TableName)
//	svcConfig := config.SvcConfig{
//		Cfg:   &config.Config{DataBase: dbcfg},
//		DbSvc: config.DbSvc{Db: dataBase},
//	}
//	dB := sqlDs{
//		sqlSvc: svcConfig.DbSvc.Db,
//		table:  svcConfig.Cfg.DataBase.TableName,
//	}
//	// table driven tests
//	tests := []struct {
//		name        string
//		tableName   string
//		dataSet     map[string]interface{}
//		dataWhere   map[string]interface{}
//		setupFunc   func()
//		cleanupFunc func()
//		filter      map[string]interface{}
//		validator   func(error)
//	}{
//		{
//			name:      "SUCCESS:: Update",
//			dataSet:   map[string]interface{}{"active": true, "active_devices": 1, "company_name": "newCompany"},
//			dataWhere: map[string]interface{}{"email": "v@mail.com", "active": false},
//			setupFunc: func() {
//				tableName := "newTemp"
//				createTestTable(t, dataBase, tableName, "user_id varchar(225) not null, email varchar(225) not null unique, company_name varchar(225), name varchar(225) not null, password varchar(225) not null DEFAULT 00000000, registered_on timestamp not null DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, updated_on timestamp not null DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, active boolean not null default false, active_devices int(50) not null default 0, salt varchar(225) not null default 0000, primary key (email) ")
//				err := dB.Insert(model.User{
//					Email:    "v@mail.com",
//					Company:  "company",
//					Password: "pass",
//					Name:     "vatsal",
//
//					RegisteredOn: time.Now(),
//				})
//				if err != nil {
//					t.Fatal(err)
//				}
//			},
//			cleanupFunc: func() {
//				tableName := "newTemp"
//				deleteTestTable(t, dataBase, tableName)
//			},
//			validator: func(err error) {
//				if err != nil {
//					t.Errorf("Want: %v, Got: %v", nil, err.Error())
//				}
//
//				user, err := dB.Get(map[string]interface{}{"email": "v@mail.com"})
//				if err != nil {
//					t.Errorf("Want: %v, Got: %v", nil, err.Error())
//				}
//				if user[0].Active != true {
//					t.Errorf("Want: %v, Got: %v", true, user[0].Active)
//				}
//			},
//		},
//		{
//			name:      "Failure:: Update",
//			dataSet:   map[string]interface{}{"active": true, "active_devices": 1},
//			dataWhere: map[string]interface{}{"email": 1},
//			setupFunc: func() {
//				tableName := "newTemp"
//				createTestTable(t, dataBase, tableName, "user_id varchar(225) not null, email varchar(225) not null unique, company_name int(225), name varchar(225) not null, password varchar(225) not null DEFAULT 00000000, registered_on timestamp not null DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, updated_on timestamp not null DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, active int(50) not null default false, active_devices boolean not null default 0, salt varchar(225) not null default 0000, primary key (email) ")
//				err := dB.Insert(model.User{
//					Email: "v@mail.com",
//
//					Password: "pass",
//					Name:     "vatsal",
//
//					RegisteredOn: time.Now(),
//				})
//				if err != nil {
//					t.Fatal(err)
//				}
//			},
//			cleanupFunc: func() {
//				tableName := "newTemp"
//				deleteTestTable(t, dataBase, tableName)
//			},
//			validator: func(err error) {
//				if err.Error() != errors.New("Error 1292: Truncated incorrect DOUBLE value: 'v@mail.com'").Error() {
//					t.Errorf("Want: %v, Got: %v", "Error 1292: Truncated incorrect DOUBLE value: 'v@mail.com'", err.Error())
//				}
//			},
//		},
//	}
//	// to execute the tests in the table
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if tt.setupFunc != nil {
//				tt.setupFunc()
//			}
//			// STEP 2: call the test function
//			err := dB.Update(tt.dataSet, tt.dataWhere)
//
//			// STEP 3: validation of output
//			if tt.validator != nil {
//				tt.validator(err)
//			}
//
//			// STEP 4: clean up/remove up all instances for the specific test case
//			if tt.cleanupFunc != nil {
//				tt.cleanupFunc()
//			}
//		})
//	}
//}
