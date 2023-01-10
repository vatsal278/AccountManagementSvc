package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/vatsal278/AccountManagmentSvc/internal/codes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	respModel "github.com/PereRohit/util/model"
	"github.com/PereRohit/util/testutil"
	"github.com/golang/mock/gomock"

	"github.com/vatsal278/AccountManagmentSvc/internal/model"
	"github.com/vatsal278/AccountManagmentSvc/pkg/mock"
)

type Reader string

func (Reader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func Test_accountManagmentSvc_HealthCheck(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name        string
		setup       func() accountManagmentSvc
		wantSvcName string
		wantMsg     string
		wantStat    bool
	}{
		{
			name: "Success",
			setup: func() accountManagmentSvc {
				mockLogic := mock.NewMockAccountManagmentSvcLogicIer(mockCtrl)
				mockLogic.EXPECT().HealthCheck().
					Return(true).Times(1)

				rec := accountManagmentSvc{
					logic: mockLogic,
				}

				return rec
			},
			wantSvcName: AccountManagmentSvcName,
			wantMsg:     "",
			wantStat:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := tt.setup()

			svcName, msg, stat := receiver.HealthCheck()

			diff := testutil.Diff(svcName, tt.wantSvcName)
			if diff != "" {
				t.Error(testutil.Callers(), diff)
			}

			diff = testutil.Diff(msg, tt.wantMsg)
			if diff != "" {
				t.Error(testutil.Callers(), diff)
			}

			diff = testutil.Diff(stat, tt.wantStat)
			if diff != "" {
				t.Error(testutil.Callers(), diff)
			}
		})
	}
}

func TestAccountManagmentSvc_CreateAccount(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name  string
		model model.NewAccount
		setup func() (*accountManagmentSvc, *http.Request)
		want  func(recorder httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			model: model.NewAccount{
				UserId: "123",
			},
			setup: func() (*accountManagmentSvc, *http.Request) {
				mockLogic := mock.NewMockAccountManagmentSvcLogicIer(mockCtrl)
				mockLogic.EXPECT().CreateAccount(gomock.Any()).Times(1).Return(&respModel.Response{
					Status:  http.StatusCreated,
					Message: codes.GetErr(codes.Success),
					Data:    nil,
				})
				svc := &accountManagmentSvc{
					logic: mockLogic,
				}
				by, err := json.Marshal(model.NewAccount{
					UserId: "123",
				})
				if err != nil {
					t.Fail()
				}
				r := httptest.NewRequest("POST", "/new_account", bytes.NewBuffer(by))
				return svc, r
			},
			want: func(rec httptest.ResponseRecorder) {
				b, err := ioutil.ReadAll(rec.Body)
				if err != nil {
					return
				}
				var response respModel.Response
				err = json.Unmarshal(b, &response)
				tempResp := &respModel.Response{
					Status:  http.StatusCreated,
					Message: codes.GetErr(codes.Success),
					Data:    nil,
				}
				if !reflect.DeepEqual(&response, tempResp) {
					t.Errorf("Want: %v, Got: %v", tempResp, &response)
				}

			},
		},
		{
			name: "Failure :: CreateAccount:: Read all failure",
			setup: func() (*accountManagmentSvc, *http.Request) {
				mockLogic := mock.NewMockAccountManagmentSvcLogicIer(mockCtrl)
				svc := &accountManagmentSvc{
					logic: mockLogic,
				}
				r := httptest.NewRequest("POST", "/register", Reader(""))
				return svc, r
			},
			want: func(rec httptest.ResponseRecorder) {
				b, err := ioutil.ReadAll(rec.Body)
				if err != nil {
					return
				}
				var response respModel.Response
				err = json.Unmarshal(b, &response)
				tempResp := &respModel.Response{
					Status:  http.StatusInternalServerError,
					Message: "request body read : test error",
					Data:    nil,
				}
				if !reflect.DeepEqual(&response, tempResp) {
					t.Errorf("Want: %v, Got: %v", tempResp, &response)
				}
			},
		},
		{
			name: "Failure :: CreateAccount:: json unmarshall failure",
			setup: func() (*accountManagmentSvc, *http.Request) {
				mockLogic := mock.NewMockAccountManagmentSvcLogicIer(mockCtrl)
				svc := &accountManagmentSvc{
					logic: mockLogic,
				}
				r := httptest.NewRequest("POST", "/register", bytes.NewBuffer([]byte("")))
				return svc, r
			},
			want: func(rec httptest.ResponseRecorder) {
				b, err := ioutil.ReadAll(rec.Body)
				if err != nil {
					return
				}
				var response respModel.Response
				err = json.Unmarshal(b, &response)
				tempResp := &respModel.Response{
					Status:  http.StatusBadRequest,
					Message: "put data into data: unexpected end of JSON input",
					Data:    nil,
				}
				if !reflect.DeepEqual(&response, tempResp) {
					t.Errorf("Want: %v, Got: %v", tempResp, &response)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			x, r := tt.setup()
			x.CreateAccount(w, r)
			tt.want(*w)
		})
	}
}
