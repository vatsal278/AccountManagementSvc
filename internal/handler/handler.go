package handler

import (
	"github.com/vatsal278/AccountManagmentSvc/internal/logic"
	"github.com/vatsal278/AccountManagmentSvc/internal/repo/datasource"
)

const AccountManagmentSvcName = "accountManagmentSvc"

//go:generate mockgen --build_flags=--mod=mod --destination=./../../pkg/mock/mock_handler.go --package=mock github.com/vatsal278/AccountManagmentSvc/internal/handler AccountManagmentSvcHandler

type AccountManagmentSvcHandler interface {
	HealthChecker
}

type accountManagmentSvc struct {
	logic logic.AccountManagmentSvcLogicIer
}

func NewAccountManagmentSvc(ds datasource.DataSourceI) AccountManagmentSvcHandler {
	svc := &accountManagmentSvc{
		logic: logic.NewAccountManagmentSvcLogic(ds),
	}
	AddHealthChecker(svc)
	return svc
}

func (svc accountManagmentSvc) HealthCheck() (svcName string, msg string, stat bool) {
	set := false
	defer func() {
		svcName = AccountManagmentSvcName
		if !set {
			msg = ""
			stat = true
		}
	}()
	stat = svc.logic.HealthCheck()
	set = true
	return
}
