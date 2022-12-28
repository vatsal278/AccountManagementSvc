package logic

import (
	"net/http"

	"github.com/PereRohit/util/log"
	respModel "github.com/PereRohit/util/model"

	"github.com/vatsal278/AccountManagmentSvc/internal/model"
	"github.com/vatsal278/AccountManagmentSvc/internal/repo/datasource"
)

//go:generate mockgen --build_flags=--mod=mod --destination=./../../pkg/mock/mock_logic.go --package=mock github.com/vatsal278/AccountManagmentSvc/internal/logic AccountManagmentSvcLogicIer

type AccountManagmentSvcLogicIer interface {
	Ping(*model.PingRequest) *respModel.Response
	HealthCheck() bool
}

type accountManagmentSvcLogic struct {
	dummyDsSvc datasource.DataSource
}

func NewAccountManagmentSvcLogic(ds datasource.DataSource) AccountManagmentSvcLogicIer {
	return &accountManagmentSvcLogic{
		dummyDsSvc: ds,
	}
}

func (l accountManagmentSvcLogic) Ping(req *model.PingRequest) *respModel.Response {
	// add business logic here
	res, err := l.dummyDsSvc.Ping(&model.PingDs{
		Data: req.Data,
	})
	if err != nil {
		log.Error("datasource error", err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: "",
			Data:    nil,
		}
	}
	return &respModel.Response{
		Status:  http.StatusOK,
		Message: "Pong",
		Data:    res,
	}
}

func (l accountManagmentSvcLogic) HealthCheck() bool {
	// check all internal services are working fine
	return l.dummyDsSvc.HealthCheck()
}
