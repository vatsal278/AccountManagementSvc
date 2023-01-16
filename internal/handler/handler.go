package handler

import (
	"github.com/PereRohit/util/log"
	"github.com/PereRohit/util/request"
	"github.com/PereRohit/util/response"
	"github.com/vatsal278/AccountManagmentSvc/internal/codes"
	"github.com/vatsal278/AccountManagmentSvc/internal/config"
	"github.com/vatsal278/AccountManagmentSvc/internal/logic"
	"github.com/vatsal278/AccountManagmentSvc/internal/model"
	jwtSvc "github.com/vatsal278/AccountManagmentSvc/internal/repo/authentication"
	"github.com/vatsal278/AccountManagmentSvc/internal/repo/datasource"
	"github.com/vatsal278/AccountManagmentSvc/pkg/session"
	"net/http"
)

const AccountManagmentSvcName = "accountManagmentSvc"

//go:generate mockgen --build_flags=--mod=mod --destination=./../../pkg/mock/mock_handler.go --package=mock github.com/vatsal278/AccountManagmentSvc/internal/handler AccountManagmentSvcHandler

type AccountManagmentSvcHandler interface {
	HealthChecker
	CreateAccount(w http.ResponseWriter, r *http.Request)
	AccountSummary(w http.ResponseWriter, r *http.Request)
}

type accountManagmentSvc struct {
	logic logic.AccountManagmentSvcLogicIer
}

func NewAccountManagmentSvc(ds datasource.DataSourceI, jwtService jwtSvc.JWTService, msgQueue config.MsgQueue, cookie config.CookieStruct) AccountManagmentSvcHandler {
	svc := &accountManagmentSvc{
		logic: logic.NewAccountManagmentSvcLogic(ds, jwtService, msgQueue, cookie),
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

func (svc accountManagmentSvc) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var newAccount model.NewAccount
	status, err := request.FromJson(r, &newAccount)
	if err != nil {
		log.Error(err)
		response.ToJson(w, status, err.Error(), nil)
		return
	}
	resp := svc.logic.CreateAccount(newAccount)
	response.ToJson(w, resp.Status, resp.Message, resp.Data)
}

func (svc accountManagmentSvc) AccountSummary(w http.ResponseWriter, r *http.Request) {
	id := session.GetSession(r.Context())
	idStr, ok := id.(string)
	if !ok {
		response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrAssertUserid), nil)
		return
	}
	resp := svc.logic.AccountDetails(idStr)
	response.ToJson(w, resp.Status, resp.Message, resp.Data)
}
