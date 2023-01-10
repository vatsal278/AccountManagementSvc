package logic

import (
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
	CreateAccount(account model.Account) *respModel.Response
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

func (l accountManagmentSvcLogic) CreateAccount(account model.Account) *respModel.Response {
	result, err := l.DsSvc.Get(map[string]interface{}{"user_id": account.Id})
	if err != nil {
		log.Error(err.Error())
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrCreatingAccount),
			Data:    nil,
		}
	}
	if len(result) != 0 {
		log.Error("acc for this id already there")
		return &respModel.Response{
			Status:  http.StatusBadRequest,
			Message: "acc for this id already there",
			Data:    nil,
		}
	}
	err = l.DsSvc.Insert(account)
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
	}(account.Id, l.msgQueue.PubId, l.msgQueue.Channel)
	return &respModel.Response{
		Status:  http.StatusCreated,
		Message: "SUCCESS",
		Data:    nil,
	}
}
