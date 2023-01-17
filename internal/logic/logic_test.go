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
func TestAccountManagmentSvcLogic_AccountSummary(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name        string
		credentials string
		setup       func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue, config.CookieStruct)
		want        func(*respModel.Response)
	}{
		{
			name:        "Success :: AccDetails",
			credentials: "123",
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue, config.CookieStruct) {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				mockJwtSvc := mock.NewMockJWTService(mockCtrl)
				var acc []model.Account
				acc = append(acc, model.Account{Id: "123", AccountNumber: 1})
				mockDs.EXPECT().Get(map[string]interface{}{"user_id": "123"}).Times(1).Return(acc, nil)
				return mockDs, mockJwtSvc, config.MsgQueue{}, config.CookieStruct{}
			},
			want: func(resp *respModel.Response) {
				var users = model.Account{Id: "123", AccountNumber: 1}
				temp := respModel.Response{
					Status:  http.StatusOK,
					Message: "SUCCESS",
					Data:    users,
				}
				if !reflect.DeepEqual(resp, &temp) {
					t.Errorf("Want: %v, Got: %v", &temp, resp)
				}
			},
		},
		{
			name:        "Failure :: AccDetails :: db err",
			credentials: "123",
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue, config.CookieStruct) {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				mockJwtSvc := mock.NewMockJWTService(mockCtrl)
				mockDs.EXPECT().Get(map[string]interface{}{"user_id": "123"}).Times(1).Return(nil, errors.New(""))
				return mockDs, mockJwtSvc, config.MsgQueue{}, config.CookieStruct{}
			},
			want: func(resp *respModel.Response) {
				temp := respModel.Response{
					Status:  http.StatusInternalServerError,
					Message: codes.GetErr(codes.ErrFetchingUser),
					Data:    nil,
				}
				if !reflect.DeepEqual(resp, &temp) {
					t.Errorf("Want: %v, Got: %v", &temp, resp)
				}
			},
		},
		{
			name:        "Failure :: AccDetails :: db err",
			credentials: "123",
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue, config.CookieStruct) {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				mockJwtSvc := mock.NewMockJWTService(mockCtrl)
				mockDs.EXPECT().Get(map[string]interface{}{"user_id": "123"}).Times(1).Return(nil, nil)
				return mockDs, mockJwtSvc, config.MsgQueue{}, config.CookieStruct{}
			},
			want: func(resp *respModel.Response) {
				temp := respModel.Response{
					Status:  http.StatusBadRequest,
					Message: codes.GetErr(codes.AccNotFound),
					Data:    nil,
				}
				if !reflect.DeepEqual(resp, &temp) {
					t.Errorf("Want: %v, Got: %v", &temp, resp)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := NewAccountManagmentSvcLogic(tt.setup())

			got := rec.AccountDetails(tt.credentials)

			tt.want(got)
		})
	}
}
func TestAccountManagmentSvcLogic_UpdateServices(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name        string
		credentials model.UpdateServices
		setup       func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue, config.CookieStruct)
		want        func(*respModel.Response)
	}{
		{
			name: "Success",
			credentials: model.UpdateServices{
				AccountNumber: 1,
				ServiceId:     "10",
				UpdateType:    "add",
			},
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue, config.CookieStruct) {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				mockDs.EXPECT().Update(map[string]interface{}{"active_services": model.ColumnUpdate{UpdateSet: "JSON_INSERT(active_services, '$.\"10\"', JSON_OBJECT())"}, "inactive_services": model.ColumnUpdate{UpdateSet: "JSON_REMOVE(inactive_services, '$.\"10\"')"}}, map[string]interface{}{"account_number": 1, "user_id": "1234"}).Times(1).Return(nil)
				return mockDs, nil, config.MsgQueue{MsgBroker: sdk.NewMsgBrokerSvc("http://localhost:9095")}, config.CookieStruct{}
			},
			want: func(resp *respModel.Response) {
				temp := respModel.Response{
					Status:  http.StatusAccepted,
					Message: "SUCCESS",
					Data:    nil,
				}
				if !reflect.DeepEqual(resp, &temp) {
					t.Errorf("Want: %v, Got: %v", temp, resp)
				}
			},
		},
		{
			name: "Success :: remove",
			credentials: model.UpdateServices{
				AccountNumber: 1,
				ServiceId:     "10",
				UpdateType:    "remove",
			},
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue, config.CookieStruct) {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				mockDs.EXPECT().Update(map[string]interface{}{"active_services": model.ColumnUpdate{UpdateSet: "JSON_REMOVE(active_services, '$.\"10\"')"}, "inactive_services": model.ColumnUpdate{UpdateSet: "JSON_INSERT(inactive_services, '$.\"10\"', JSON_OBJECT())"}}, map[string]interface{}{"account_number": 1, "user_id": "1234"}).Times(1).Return(nil)
				return mockDs, nil, config.MsgQueue{MsgBroker: sdk.NewMsgBrokerSvc("http://localhost:9095")}, config.CookieStruct{}
			},
			want: func(resp *respModel.Response) {
				temp := respModel.Response{
					Status:  http.StatusAccepted,
					Message: "SUCCESS",
					Data:    nil,
				}
				if !reflect.DeepEqual(resp, &temp) {
					t.Errorf("Want: %v, Got: %v", temp, resp)
				}
			},
		},
		{
			name: "Failure::DB ERR",
			credentials: model.UpdateServices{
				AccountNumber: 1,
				ServiceId:     "10",
				UpdateType:    "add",
			},
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue, config.CookieStruct) {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				mockDs.EXPECT().Update(gomock.Any(), gomock.Any()).Times(1).Return(errors.New("DB ERR"))
				return mockDs, nil, config.MsgQueue{MsgBroker: sdk.NewMsgBrokerSvc("http://localhost:9095")}, config.CookieStruct{}
			},
			want: func(resp *respModel.Response) {
				temp := respModel.Response{
					Status:  http.StatusInternalServerError,
					Message: codes.GetErr(codes.ErrUpdatingServices),
					Data:    nil,
				}
				if !reflect.DeepEqual(resp, &temp) {
					t.Errorf("Want: %v, Got: %v", temp, resp)
				}
			},
		},
		{
			name: "Failure::switch default case",
			credentials: model.UpdateServices{
				AccountNumber: 1,
				ServiceId:     "10",
				UpdateType:    "",
			},
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue, config.CookieStruct) {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				return mockDs, nil, config.MsgQueue{MsgBroker: sdk.NewMsgBrokerSvc("http://localhost:9095")}, config.CookieStruct{}
			},
			want: func(resp *respModel.Response) {
				temp := respModel.Response{
					Status:  http.StatusBadRequest,
					Message: codes.GetErr(codes.ErrUpdatingServices),
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

			got := rec.UpdateServices("1234", tt.credentials)

			tt.want(got)
		})
	}
}
func TestAccountManagmentSvcLogic_UpdateTransaction(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name        string
		credentials model.UpdateTransaction
		setup       func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue, config.CookieStruct)
		want        func(*respModel.Response)
	}{
		{
			name: "Success::DEBIT",
			credentials: model.UpdateTransaction{
				AccountNumber:   1,
				Amount:          "1000",
				TransactionType: "debit",
			},
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue, config.CookieStruct) {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				mockDs.EXPECT().Update(map[string]interface{}{"spends": model.ColumnUpdate{UpdateSet: "spends + 1000"}}, map[string]interface{}{"account_number": 1}).Times(1).Return(nil)
				return mockDs, nil, config.MsgQueue{MsgBroker: sdk.NewMsgBrokerSvc("http://localhost:9095")}, config.CookieStruct{}
			},
			want: func(resp *respModel.Response) {
				temp := respModel.Response{
					Status:  http.StatusAccepted,
					Message: "SUCCESS",
					Data:    nil,
				}
				if !reflect.DeepEqual(resp, &temp) {
					t.Errorf("Want: %v, Got: %v", temp, resp)
				}
			},
		},
		{
			name: "Success :: CREDIT",
			credentials: model.UpdateTransaction{
				AccountNumber:   1,
				Amount:          "1000",
				TransactionType: "credit",
			},
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue, config.CookieStruct) {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				mockDs.EXPECT().Update(map[string]interface{}{"income": model.ColumnUpdate{UpdateSet: "income + 1000"}}, map[string]interface{}{"account_number": 1}).Times(1).Return(nil)
				return mockDs, nil, config.MsgQueue{MsgBroker: sdk.NewMsgBrokerSvc("http://localhost:9095")}, config.CookieStruct{}
			},
			want: func(resp *respModel.Response) {
				temp := respModel.Response{
					Status:  http.StatusAccepted,
					Message: "SUCCESS",
					Data:    nil,
				}
				if !reflect.DeepEqual(resp, &temp) {
					t.Errorf("Want: %v, Got: %v", temp, resp)
				}
			},
		},
		{
			name: "Failure::DB ERR",
			credentials: model.UpdateTransaction{
				AccountNumber:   1,
				Amount:          "1000",
				TransactionType: "debit",
			},
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue, config.CookieStruct) {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				mockDs.EXPECT().Update(gomock.Any(), gomock.Any()).Times(1).Return(errors.New("DB ERR"))
				return mockDs, nil, config.MsgQueue{MsgBroker: sdk.NewMsgBrokerSvc("http://localhost:9095")}, config.CookieStruct{}
			},
			want: func(resp *respModel.Response) {
				temp := respModel.Response{
					Status:  http.StatusInternalServerError,
					Message: codes.GetErr(codes.ErrUpdatingTransaction),
					Data:    nil,
				}
				if !reflect.DeepEqual(resp, &temp) {
					t.Errorf("Want: %v, Got: %v", temp, resp)
				}
			},
		},
		{
			name: "Failure::default switch case",
			credentials: model.UpdateTransaction{
				AccountNumber:   1,
				Amount:          "1000",
				TransactionType: "",
			},
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue, config.CookieStruct) {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				return mockDs, nil, config.MsgQueue{MsgBroker: sdk.NewMsgBrokerSvc("http://localhost:9095")}, config.CookieStruct{}
			},
			want: func(resp *respModel.Response) {
				temp := respModel.Response{
					Status:  http.StatusBadRequest,
					Message: codes.GetErr(codes.ErrUpdatingTransaction),
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

			got := rec.UpdateTransaction(tt.credentials)

			tt.want(got)
		})
	}
}
