package logic

import (
	"github.com/vatsal278/AccountManagmentSvc/internal/repo/datasource"
)

//go:generate mockgen --build_flags=--mod=mod --destination=./../../pkg/mock/mock_logic.go --package=mock github.com/vatsal278/AccountManagmentSvc/internal/logic AccountManagmentSvcLogicIer

type AccountManagmentSvcLogicIer interface {
	HealthCheck() bool
}

type accountManagmentSvcLogic struct {
	DsSvc datasource.DataSourceI
}

func NewAccountManagmentSvcLogic(ds datasource.DataSourceI) AccountManagmentSvcLogicIer {
	return &accountManagmentSvcLogic{
		DsSvc: ds,
	}
}

func (l accountManagmentSvcLogic) HealthCheck() bool {
	// check all internal services are working fine
	return l.DsSvc.HealthCheck()
}
