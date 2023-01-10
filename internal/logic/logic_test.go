package logic

import (
	"errors"
	respModel "github.com/PereRohit/util/model"
	"github.com/PereRohit/util/testutil"
	"github.com/golang/mock/gomock"
	"github.com/vatsal278/AccountManagmentSvc/internal/codes"
	"github.com/vatsal278/AccountManagmentSvc/internal/config"
	"github.com/vatsal278/AccountManagmentSvc/internal/repo/authentication"
	"github.com/vatsal278/msgbroker/pkg/sdk"
	"net/http"
	"reflect"
	"testing"

	"github.com/vatsal278/AccountManagmentSvc/internal/model"
	"github.com/vatsal278/AccountManagmentSvc/internal/repo/datasource"
	"github.com/vatsal278/AccountManagmentSvc/pkg/mock"
)

func TestAccountManagmentSvcLogic_HealthCheck(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name  string
		setup func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue, config.CookieStruct)
		want  bool
	}{
		{
			name: "Success",
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue, config.CookieStruct) {
				mockDs := mock.NewMockDataSourceI(mockCtrl)

				mockDs.EXPECT().HealthCheck().Times(1).
					Return(true)

				return mockDs, nil, config.MsgQueue{}, config.CookieStruct{}
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := NewAccountManagmentSvcLogic(tt.setup())

			got := rec.HealthCheck()

			diff := testutil.Diff(got, tt.want)
			if diff != "" {
				t.Error(testutil.Callers(), diff)
			}
		})
	}
}

func TestAccountManagmentSvcLogic_CreateAccount(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name        string
		credentials model.NewAccount
		setup       func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue, config.CookieStruct)
		want        func(*respModel.Response)
	}{
		{
			name: "Success",
			credentials: model.NewAccount{
				UserId: "123",
			},
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue, config.CookieStruct) {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				mockDs.EXPECT().Get(map[string]interface{}{"user_id": "123"}).Times(1).Return([]model.Account{}, nil)
				mockDs.EXPECT().Insert(gomock.Any()).Times(1).Return(nil)
				return mockDs, nil, config.MsgQueue{MsgBroker: sdk.NewMsgBrokerSvc("http://localhost:9095")}, config.CookieStruct{}
			},
			want: func(resp *respModel.Response) {
				temp := respModel.Response{
					Status:  http.StatusCreated,
					Message: "SUCCESS",
					Data:    nil,
				}
				if !reflect.DeepEqual(resp, &temp) {
					t.Errorf("Want: %v, Got: %v", temp, resp)
				}
			},
		},
		{
			name: "Success :: Push msg failure",
			credentials: model.NewAccount{
				UserId: "123",
			},
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue, config.CookieStruct) {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				mockDs.EXPECT().Get(map[string]interface{}{"user_id": "123"}).Times(1).Return([]model.Account{}, nil)
				mockDs.EXPECT().Insert(gomock.Any()).Times(1).Return(nil)
				return mockDs, nil, config.MsgQueue{MsgBroker: sdk.NewMsgBrokerSvc("")}, config.CookieStruct{}
			},
			want: func(resp *respModel.Response) {
				temp := respModel.Response{
					Status:  http.StatusCreated,
					Message: "SUCCESS",
					Data:    nil,
				}
				if !reflect.DeepEqual(resp, &temp) {
					t.Errorf("Want: %v, Got: %v", temp, resp)
				}
			},
		},
		{
			name: "Failure::Get from db err",
			credentials: model.NewAccount{
				UserId: "123",
			},
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue, config.CookieStruct) {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				mockDs.EXPECT().Get(gomock.Any()).Times(1).Return(nil, errors.New(""))
				return mockDs, nil, config.MsgQueue{MsgBroker: sdk.NewMsgBrokerSvc("http://localhost:9091")}, config.CookieStruct{}
			},
			want: func(resp *respModel.Response) {
				temp := respModel.Response{
					Status:  http.StatusInternalServerError,
					Message: codes.GetErr(codes.ErrCreatingAccount),
					Data:    nil,
				}
				if !reflect.DeepEqual(resp, &temp) {
					t.Errorf("Want: %v, Got: %v", temp, resp)
				}
			},
		},
		{
			name: "Failure:: Email already exists",
			credentials: model.NewAccount{
				UserId: "123",
			},
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue, config.CookieStruct) {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				var users []model.Account
				users = append(users, model.Account{Id: "123"})
				mockDs.EXPECT().Get(gomock.Any()).Times(1).Return(users, nil)
				return mockDs, nil, config.MsgQueue{MsgBroker: sdk.NewMsgBrokerSvc("http://localhost:9091")}, config.CookieStruct{}
			},
			want: func(resp *respModel.Response) {
				temp := respModel.Response{
					Status:  http.StatusBadRequest,
					Message: codes.GetErr(codes.ErrAccExists),
					Data:    nil,
				}
				if !reflect.DeepEqual(resp, &temp) {
					t.Errorf("Want: %v, Got: %v", temp, resp)
				}
			},
		},
		{
			name: "Failure:: Error Inserting in db",
			credentials: model.NewAccount{
				UserId: "123",
			},
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue, config.CookieStruct) {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				mockDs.EXPECT().Get(map[string]interface{}{"user_id": "123"}).Times(1).Return([]model.Account{}, nil)
				mockDs.EXPECT().Insert(gomock.Any()).Times(1).Return(errors.New(""))
				return mockDs, nil, config.MsgQueue{MsgBroker: sdk.NewMsgBrokerSvc("http://localhost:9091")}, config.CookieStruct{}
			},
			want: func(resp *respModel.Response) {
				temp := respModel.Response{
					Status:  http.StatusBadRequest,
					Message: codes.GetErr(codes.ErrCreatingAccount),
					Data:    nil,
				}
				if !reflect.DeepEqual(resp, &temp) {
					t.Errorf("Want: %v, Got: %v", temp, resp)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := NewAccountManagmentSvcLogic(tt.setup())

			got := rec.CreateAccount(tt.credentials)

			tt.want(got)
		})
	}
}
