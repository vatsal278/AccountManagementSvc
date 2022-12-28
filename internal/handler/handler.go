package handler

import (
	"fmt"
	"net/http"

	"github.com/PereRohit/util/request"
	"github.com/PereRohit/util/response"

	"github.com/vatsal278/AccountManagmentSvc/internal/logic"
	"github.com/vatsal278/AccountManagmentSvc/internal/model"
	"github.com/vatsal278/AccountManagmentSvc/internal/repo/datasource"
)

const AccountManagmentSvcName = "accountManagmentSvc"

//go:generate mockgen --build_flags=--mod=mod --destination=./../../pkg/mock/mock_handler.go --package=mock github.com/vatsal278/AccountManagmentSvc/internal/handler AccountManagmentSvcHandler

type AccountManagmentSvcHandler interface {
	HealthChecker
	Ping(w http.ResponseWriter, r *http.Request)
}

type accountManagmentSvc struct {
	logic logic.AccountManagmentSvcLogicIer
}

func NewAccountManagmentSvc(ds datasource.DataSource) AccountManagmentSvcHandler {
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

func (svc accountManagmentSvc) Ping(w http.ResponseWriter, r *http.Request) {
	req := &model.PingRequest{}

	suggestedCode, err := request.FromJson(r, req)
	if err != nil {
		response.ToJson(w, suggestedCode, fmt.Sprintf("FAILED: %s", err.Error()), nil)
		return
	}
	// call logic
	resp := svc.logic.Ping(req)
	response.ToJson(w, resp.Status, resp.Message, resp.Data)
	return
}
