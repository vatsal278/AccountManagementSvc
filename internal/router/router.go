package router

import (
	"net/http"

	"github.com/PereRohit/util/constant"
	"github.com/PereRohit/util/middleware"
	"github.com/gorilla/mux"

	"github.com/vatsal278/AccountManagmentSvc/internal/config"
	"github.com/vatsal278/AccountManagmentSvc/internal/handler"
	"github.com/vatsal278/AccountManagmentSvc/internal/repo/datasource"
)

func Register(svcCfg *config.SvcConfig) *mux.Router {
	m := mux.NewRouter()

	// group all routes for specific version. e.g.: /v1
	if svcCfg.ServiceRouteVersion != "" {
		m = m.PathPrefix("/" + svcCfg.ServiceRouteVersion).Subrouter()
	}

	m.StrictSlash(true)
	m.Use(middleware.RequestHijacker)
	m.Use(middleware.RecoverPanic)

	commons := handler.NewCommonSvc()
	m.HandleFunc(constant.HealthRoute, commons.HealthCheck).Methods(http.MethodGet)
	m.NotFoundHandler = http.HandlerFunc(commons.RouteNotFound)
	m.MethodNotAllowedHandler = http.HandlerFunc(commons.MethodNotAllowed)

	// attach routes for services below
	m = attachAccountManagmentSvcRoutes(m, svcCfg)

	return m
}

func attachAccountManagmentSvcRoutes(m *mux.Router, svcCfg *config.SvcConfig) *mux.Router {
	dataSource := datasource.NewDummyDs(&svcCfg.DummySvc)

	svc := handler.NewAccountManagmentSvc(dataSource)

	m.HandleFunc("/ping", svc.Ping).Methods(http.MethodPost)

	return m
}
