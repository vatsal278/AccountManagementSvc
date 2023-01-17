package logic

import (
	"errors"
	"fmt"
	"github.com/PereRohit/util/log"
	respModel "github.com/PereRohit/util/model"
	"github.com/vatsal278/AccountManagmentSvc/internal/codes"
	"github.com/vatsal278/AccountManagmentSvc/internal/config"
	"github.com/vatsal278/AccountManagmentSvc/internal/model"
	jwtSvc "github.com/vatsal278/AccountManagmentSvc/internal/repo/authentication"
	"github.com/vatsal278/AccountManagmentSvc/internal/repo/datasource"
	"net/http"
)

//go:generate mockgen --build_flags=--mod=mod --destination=./../../pkg/mock/mock_logic.go --package=mock github.com/vatsal278/AccountManagmentSvc/internal/logic AccountManagmentSvcLogicIer

type AccountManagmentSvcLogicIer interface {
	HealthCheck() bool
	CreateAccount(account model.NewAccount) *respModel.Response
	AccountDetails(id string) *respModel.Response
	UpdateServices(id string, services model.UpdateServices) *respModel.Response
	UpdateTransaction(transaction model.UpdateTransaction) *respModel.Response
}

type accountManagmentSvcLogic struct {
	DsSvc      datasource.DataSourceI
	jwtService jwtSvc.JWTService
	msgQueue   config.MsgQueue
	cookie     config.CookieStruct
}

func NewAccountManagmentSvcLogic(ds datasource.DataSourceI, jwtService jwtSvc.JWTService, msgQueue config.MsgQueue, cookie config.CookieStruct) AccountManagmentSvcLogicIer {
	return &accountManagmentSvcLogic{
		DsSvc:      ds,
		jwtService: jwtService,
		msgQueue:   msgQueue,
		cookie:     cookie,
	}
}

func (l accountManagmentSvcLogic) HealthCheck() bool {
	// check all internal services are working fine
	return l.DsSvc.HealthCheck()
}

func (l accountManagmentSvcLogic) CreateAccount(account model.NewAccount) *respModel.Response {
	result, err := l.DsSvc.Get(map[string]interface{}{"user_id": account.UserId})
	if err != nil {
		log.Error(err.Error())
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrCreatingAccount),
			Data:    nil,
		}
	}
	if len(result) != 0 {
		log.Error(codes.GetErr(codes.ErrAccExists))
		return &respModel.Response{
			Status:  http.StatusBadRequest,
			Message: codes.GetErr(codes.ErrAccExists),
			Data:    nil,
		}
	}
	err = l.DsSvc.Insert(model.Account{Id: account.UserId})
	if err != nil {
		log.Error(codes.GetErr(codes.ErrCreatingAccount))
		return &respModel.Response{
			Status:  http.StatusBadRequest,
			Message: codes.GetErr(codes.ErrCreatingAccount),
			Data:    nil,
		}
	}
	go func(userId string, pubId string, channel string) {
		userID := fmt.Sprintf(`{"user_id":"%s"}`, userId)
		err := l.msgQueue.MsgBroker.PushMsg(userID, pubId, channel)
		if err != nil {
			log.Error(err)
			return
		}
		return
	}(account.UserId, l.msgQueue.PubId, l.msgQueue.Channel)
	return &respModel.Response{
		Status:  http.StatusCreated,
		Message: "SUCCESS",
		Data:    nil,
	}
}

func (l accountManagmentSvcLogic) AccountDetails(id string) *respModel.Response {
	acc, err := l.DsSvc.Get(map[string]interface{}{"user_id": id})
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrFetchingUser),
			Data:    nil,
		}
	}
	if len(acc) == 0 {
		return &respModel.Response{
			Status:  http.StatusBadRequest,
			Message: codes.GetErr(codes.AccNotFound),
			Data:    nil,
		}
	}
	resp := model.AccountSummary{
		AccountNumber:    acc[0].AccountNumber,
		Income:           acc[0].Income,
		Spends:           acc[0].Spends,
		ActiveServices:   acc[0].ActiveServices,
		InactiveServices: acc[0].InactiveServices,
	}
	return &respModel.Response{
		Status:  http.StatusOK,
		Message: "SUCCESS",
		Data:    resp,
	}
}

func (l accountManagmentSvcLogic) UpdateServices(id string, services model.UpdateServices) *respModel.Response {
	var query map[string]interface{}
	switch services.UpdateType {
	case "add":
		insertQuery := fmt.Sprintf("JSON_INSERT(%s, '$.\"%s\"', JSON_OBJECT())", "active_services", services.ServiceId)
		removeQuery := fmt.Sprintf("JSON_REMOVE(%s, '$.\"%s\"')", "inactive_services", services.ServiceId)
		query = map[string]interface{}{"active_services": model.ColumnUpdate{UpdateSet: insertQuery}, "inactive_services": model.ColumnUpdate{UpdateSet: removeQuery}}
	case "remove":
		insertQuery := fmt.Sprintf("JSON_INSERT(%s, '$.\"%s\"', JSON_OBJECT())", "inactive_services", services.ServiceId)
		removeQuery := fmt.Sprintf("JSON_REMOVE(%s, '$.\"%s\"')", "active_services", services.ServiceId)
		query = map[string]interface{}{"inactive_services": model.ColumnUpdate{UpdateSet: insertQuery}, "active_services": model.ColumnUpdate{UpdateSet: removeQuery}}
	default:
		log.Error(errors.New("incorrect update query "))
		return &respModel.Response{
			Status:  http.StatusBadRequest,
			Message: codes.GetErr(codes.ErrUpdatingServices),
			Data:    nil,
		}
	}
	err := l.DsSvc.Update(query, map[string]interface{}{"user_id": id, "account_number": services.AccountNumber})
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrUpdatingServices),
			Data:    nil,
		}
	}
	return &respModel.Response{
		Status:  http.StatusAccepted,
		Message: "SUCCESS",
		Data:    nil,
	}
}
func (l accountManagmentSvcLogic) UpdateTransaction(transaction model.UpdateTransaction) *respModel.Response {
	var query map[string]interface{}
	switch transaction.TransactionType {
	case "debit":
		spends := fmt.Sprintf("spends + %s", transaction.Amount)
		query = map[string]interface{}{"spends": model.ColumnUpdate{UpdateSet: spends}}
	case "credit":
		spends := fmt.Sprintf("income + %s", transaction.Amount)
		query = map[string]interface{}{"income": model.ColumnUpdate{UpdateSet: spends}}
	default:
		log.Error(errors.New("incorrect transaction type "))
		return &respModel.Response{
			Status:  http.StatusBadRequest,
			Message: codes.GetErr(codes.ErrUpdatingServices),
			Data:    nil,
		}
	}
	err := l.DsSvc.Update(query, map[string]interface{}{"account_number": transaction.AccountNumber})
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrUpdatingTransaction),
			Data:    nil,
		}
	}
	return &respModel.Response{
		Status:  http.StatusAccepted,
		Message: "SUCCESS",
		Data:    nil,
	}
}
