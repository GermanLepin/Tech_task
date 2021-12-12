package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"tech_task/pkg/service"
	mock_service "tech_task/pkg/service/mocks"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"gotest.tools/assert"
)

func TestHandler_BalanceInfo(t *testing.T) {
	type mockBehavior func(r *mock_service.MockBalanceInfo, id int64)
	ctx := context.Background()

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            int64
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"id":"1"}`,
			inputUser: 1,
			mockBehavior: func(r *mock_service.MockBalanceInfo, id int64) {
				var uid int64 = 1
				var balance float64 = 830.55
				var err error = nil
				r.EXPECT().BalanceInfoUser(ctx, id).Return(uid, balance, err)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "{\"user id\":1,\"balance\":830.55}\n",
		},
		{
			name:                 "Wrong Input",
			inputBody:            `{"id":"-1"}`,
			inputUser:            -1,
			mockBehavior:         func(r *mock_service.MockBalanceInfo, id int64) {},
			expectedStatusCode:   400,
			expectedResponseBody: "{\"error\":\"incorrect value id user\"}\n",
		},
		{
			name:      "User not found",
			inputBody: `{"id":"99999999"}`,
			inputUser: 99999999,
			mockBehavior: func(r *mock_service.MockBalanceInfo, id int64) {
				var uid int64 = 0
				var balance float64 = 0
				var err error = errors.New("{\"error\":\"User not found\"}\n")
				r.EXPECT().BalanceInfoUser(ctx, id).Return(uid, balance, err)
			},
			expectedStatusCode:   400,
			expectedResponseBody: "{\"error\":\"User not found\"}\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockBalanceInfoUser(c)
			test.mockBehavior(repo, test.inputUser)

			services := &service.Service{BalanceInfo: repo}
			handler := Handler{services}

			r := chi.NewRouter()
			r.Get("/balance-info", handler.BalanceInfo)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/balance-info",
				bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}