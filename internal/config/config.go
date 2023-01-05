package config

import (
	"database/sql"
	"fmt"
	"github.com/PereRohit/util/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/vatsal278/AccountManagmentSvc/internal/model"
	"log"
)

type Config struct {
	ServiceRouteVersion string              `json:"service_route_version"`
	ServerConfig        config.ServerConfig `json:"server_config"`
	// add custom config structs below for any internal services
	DataBase DbCfg `json:"db_svc"`
}

type SvcConfig struct {
	Cfg                 *Config
	ServiceRouteVersion string
	SvrCfg              config.ServerConfig
	// add internal services after init
	DbSvc DbSvc
}

type DbSvc struct {
	Db *sql.DB
}
type DbCfg struct {
	Port      string `json:"dbPort"`
	Host      string `json:"dbHost"`
	Driver    string `json:"dbDriver"`
	User      string `json:"dbUser"`
	Pass      string `json:"dbPass"`
	DbName    string `json:"dbName"`
	TableName string `json:"tableName"`
}
type JWTSvc struct {
	JwtSvc jwtSvc.JWTService
}

func Connect(cfg DbCfg, tableName string) *sql.DB {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True", cfg.User, cfg.Pass, cfg.Host, cfg.Port)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		panic(err.Error())
	}
	dbString := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s ;", cfg.DbName)
	prepare, err := db.Prepare(dbString)
	if err != nil {
		log.Print(err)
		return nil
	}
	_, err = prepare.Exec()
	if err != nil {
		log.Print(err)
		return nil
	}
	db.Close()
	connectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.DbName)
	db, err = sql.Open(cfg.Driver, connectionString)
	if err != nil {
		panic(err.Error())
	}
	x := fmt.Sprintf("create table if not exists %s", tableName)
	_, err = db.Exec(x + model.Schema)
	if err != nil {
		log.Fatal(err.Error())
	}
	return db
}

func InitSvcConfig(cfg Config) *SvcConfig {
	// init required services and assign to the service struct fields
	dataBase := Connect(cfg.DataBase, cfg.DataBase.TableName)
	return &SvcConfig{
		Cfg:                 &cfg,
		ServiceRouteVersion: cfg.ServiceRouteVersion,
		SvrCfg:              cfg.ServerConfig,
		DbSvc:               DbSvc{Db: dataBase},
	}
}
