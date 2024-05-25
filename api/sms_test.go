package api

import (
	"bytes"
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSmsGet(t *testing.T) {
	newRequest := func(smsId string, modemId string, logger *slog.Logger) *http.Request {
		r := httptest.NewRequest("GET", "http://127.0.0.1:8743/sms/", nil)

		ctx := context.WithValue(r.Context(), "logger", logger)
		ctx = context.WithValue(ctx, "modem", modemId)

		if len(smsId) > 0 {
			routeContext := chi.NewRouteContext()
			routeContext.URLParams.Add("id", smsId)
			ctx = context.WithValue(ctx, chi.RouteCtxKey, routeContext)
		}

		return r.WithContext(ctx)
	}

	tests := []struct {
		name               string
		modemId            string
		smsId              string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			"List success",
			"1",
			"",
			http.StatusOK,
			`{"modem.messaging.sms":["/org/freedesktop/ModemManager1/SMS/0"]}`,
		},
		{
			"List error",
			"1",
			"",
			http.StatusInternalServerError,
			`{"error":"internal server error"}`,
		},
		{
			"Single success",
			"",
			"/org/freedesktop/ModemManager1/SMS/0",
			http.StatusOK,
			`{
  "sms": {
    "content": {
      "data": "--",
      "number": "+43123456789",
      "text": "Ping"
    },
    "dbus-path": "/org/freedesktop/ModemManager1/SMS/0",
    "properties": {
      "class": "--",
      "delivery-report": "not requested",
      "delivery-state": "--",
      "discharge-timestamp": "--",
      "message-reference": "6",
      "pdu-type": "submit",
      "service-category": "--",
      "smsc": "--",
      "state": "sent",
      "storage": "--",
      "teleservice-id": "--",
      "timestamp": "--",
      "validity": "--"
    }
  }
}`,
		},
		{
			"Single error",
			"",
			"/org/freedesktop/ModemManager1/SMS/0",
			http.StatusInternalServerError,
			`{"error":"internal server error"}`,
		},
		{
			"Single invalid sms id",
			"",
			"#ß?",
			http.StatusBadRequest,
			`{"error":"invalid ID"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api, backendMock := ProvideTestApi(t)

			var buff bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&buff, nil))

			r := newRequest(tt.smsId, tt.modemId, logger)
			w := httptest.NewRecorder()

			if len(tt.smsId) > 0 {
				if tt.expectedStatusCode == http.StatusInternalServerError {
					backendMock.EXPECT().Exec("-s", tt.smsId).Return(nil, errors.New("failed"))
				} else if tt.expectedStatusCode == http.StatusBadRequest {
				} else {
					backendMock.EXPECT().Exec("-s", tt.smsId).Return([]byte(tt.expectedBody), nil)
				}
			} else {
				if tt.expectedStatusCode == http.StatusInternalServerError {
					backendMock.EXPECT().ExecModem(tt.modemId, "--messaging-list-sms").Return(nil, errors.New("failed"))
				} else {
					backendMock.EXPECT().ExecModem(tt.modemId, "--messaging-list-sms").Return([]byte(tt.expectedBody), nil)
				}
			}
			api.SmsGet(w, r)
			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)
			assert.Equal(t, tt.expectedBody, strings.Trim(string(body), " \n"))
		})
	}
}

func TestSmsSend(t *testing.T) {
	newRequest := func(body string, modemId string, logger *slog.Logger) *http.Request {
		var bodyReader io.Reader
		if len(body) > 0 {
			bodyReader = strings.NewReader(body)
		}

		r := httptest.NewRequest("POST", "http://127.0.0.1:8743/sms", bodyReader)

		ctx := context.WithValue(r.Context(), "logger", logger)
		ctx = context.WithValue(ctx, "modem", modemId)

		return r.WithContext(ctx)
	}

	tests := []struct {
		name               string
		modemId            string
		requestBody        string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			"Success",
			"",
			`{"number":"+436641234567","text":"Ping"}`,
			http.StatusOK,
			`{"message":"successfully sent the SMS"}`,
		},
		{
			"Invalid json",
			"",
			`{"foo":"+`,
			http.StatusBadRequest,
			`{"error":"unexpected EOF"}`,
		},
		{
			"Invalid request data",
			"",
			`{"foo":"+436641234567","xxx":"Ping"}`,
			http.StatusBadRequest,
			`{"error":"Key: 'SmsRequestData.Number' Error:Field validation for 'Number' failed on the 'required' tag\nKey: 'SmsRequestData.Text' Error:Field validation for 'Text' failed on the 'required' tag"}`,
		},
		{
			"Create SMS failed",
			"",
			`{"number":"+436641234567","text":"Ping"}`,
			http.StatusInternalServerError,
			`{"error":"internal server error"}`,
		},
		{
			"Create SMS invalid json",
			"",
			`{"number":"+436641234567","text":"Ping"}`,
			http.StatusInternalServerError,
			`{"error":"internal server error"}`,
		},
		{
			"Send SMS failed",
			"",
			`{"number":"+436641234567","text":"Ping"}`,
			http.StatusInternalServerError,
			`{"error":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api, backendMock := ProvideTestApi(t)

			var buff bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&buff, nil))

			r := newRequest(tt.requestBody, tt.modemId, logger)
			w := httptest.NewRecorder()

			if tt.name == "Success" {
				backendMock.EXPECT().ExecModem(tt.modemId, "--messaging-create-sms=number='+436641234567',text='Ping'").Return([]byte(`{"modem":{"messaging":{"created-sms":"/org/freedesktop/ModemManager1/SMS/0"}}}`), nil)
				backendMock.EXPECT().Exec("-s", "/org/freedesktop/ModemManager1/SMS/0", "--send").Return([]byte(`successfully sent the SMS`), nil)
			} else if tt.name == "Create SMS failed" {
				backendMock.EXPECT().ExecModem(tt.modemId, "--messaging-create-sms=number='+436641234567',text='Ping'").Return([]byte(`create SMS failed`), errors.New("failed"))
			} else if tt.name == "Create SMS invalid json" {
				backendMock.EXPECT().ExecModem(tt.modemId, "--messaging-create-sms=number='+436641234567',text='Ping'").Return([]byte(`{"m`), nil)
			} else if tt.name == "Send SMS failed" {
				backendMock.EXPECT().ExecModem(tt.modemId, "--messaging-create-sms=number='+436641234567',text='Ping'").Return([]byte(`{"modem":{"messaging":{"created-sms":"/org/freedesktop/ModemManager1/SMS/0"}}}`), nil)
				backendMock.EXPECT().Exec("-s", "/org/freedesktop/ModemManager1/SMS/0", "--send").Return([]byte(`failed`), errors.New("failed"))
			}

			api.SmsSend(w, r)
			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)
			assert.Equal(t, tt.expectedBody, strings.Trim(string(body), " \n"))
		})
	}
}

func TestSmsDelete(t *testing.T) {
	newRequest := func(smsId string, modemId string, logger *slog.Logger) *http.Request {
		r := httptest.NewRequest("DELETE", "http://127.0.0.1:8743/sms/", nil)

		ctx := context.WithValue(r.Context(), "logger", logger)
		ctx = context.WithValue(ctx, "modem", modemId)

		if len(smsId) > 0 {
			routeContext := chi.NewRouteContext()
			routeContext.URLParams.Add("id", smsId)
			ctx = context.WithValue(ctx, chi.RouteCtxKey, routeContext)
		}

		return r.WithContext(ctx)
	}

	tests := []struct {
		name               string
		modemId            string
		smsId              string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			"Success",
			"1",
			"/org/freedesktop/ModemManager1/SMS/0",
			http.StatusOK,
			`{"message":"successfully deleted SMS from modem"}`,
		},
		{
			"Invalid SMS id",
			"",
			"&#'",
			http.StatusBadRequest,
			`{"error":"invalid ID"}`,
		},
		{
			"Delete failed",
			"",
			"12",
			http.StatusInternalServerError,
			`{"error":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api, backendMock := ProvideTestApi(t)

			var buff bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&buff, nil))

			r := newRequest(tt.smsId, tt.modemId, logger)
			w := httptest.NewRecorder()

			if tt.name == "Success" {
				backendMock.EXPECT().ExecModem(tt.modemId, "--messaging-delete-sms="+tt.smsId).Return([]byte("successfully deleted SMS from modem"), nil)
			} else if tt.name == "Delete failed" {
				backendMock.EXPECT().ExecModem(tt.modemId, "--messaging-delete-sms="+tt.smsId).Return([]byte("failed"), errors.New("failed"))
			}

			api.SmsDelete(w, r)
			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)
			assert.Equal(t, tt.expectedBody, strings.Trim(string(body), " \n"))
		})
	}
}
