package http

import (
	"bytes"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testInvalidJson string

func (b testInvalidJson) MarshalJSON() ([]byte, error) {
	return nil, errors.New("encoding error")
}

func TestWriteJson(t *testing.T) {

	type car struct {
		Brand string
		Year  int
	}

	type encodingError struct {
		Brand testInvalidJson
		Year  int
	}

	tests := []struct {
		name   string
		status int
		data   any
		expect string
	}{
		{
			"Struct",
			http.StatusOK,
			car{"Toyota", 2005},
			`{"Brand":"Toyota","Year":2005}`,
		},
		{
			"Integer",
			http.StatusOK,
			12345,
			"12345",
		},
		{
			"String",
			http.StatusOK,
			"hello",
			`"hello"`,
		},
		{
			"EncodingError",
			http.StatusInternalServerError,
			encodingError{"Tesla", 2023},
			``,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			WriteJson(rec, test.status, test.data)

			resp := rec.Result()
			body, _ := io.ReadAll(resp.Body)

			assert.Equal(t, test.status, resp.StatusCode)
			assert.Equal(t, test.expect, strings.Trim(string(body), " \n"))
		})
	}
}

func TestWriteStatus(t *testing.T) {
	tests := []struct {
		name   string
		status int
	}{
		{
			"ok",
			http.StatusOK,
		},
		{
			"NotFound",
			http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			WriteStatus(rec, test.status)

			resp := rec.Result()

			assert.Equal(t, test.status, resp.StatusCode)
			assert.Equal(t, resp.Header.Get("Content-Type"), "application/json")
		})
	}
}

func TestWriteError(t *testing.T) {
	tests := []struct {
		name               string
		err                error
		statusCode         int
		expectedStatusCode int
		expectedBody       string
	}{
		{
			"unspecified error",
			nil,
			0,
			http.StatusBadRequest,
			`{"error":"unspecified error"}`,
		},
		{
			"bad request",
			errors.New("bad request"),
			http.StatusBadRequest,
			http.StatusBadRequest,
			`{"error":"bad request"}`,
		},
		{
			"specific status and error",
			errors.New("not found"),
			http.StatusNotFound,
			http.StatusNotFound,
			`{"error":"not found"}`,
		},
		{
			"internal server error",
			errors.New("something weird happened"),
			http.StatusInternalServerError,
			http.StatusInternalServerError,
			`{"error":"internal server error"}`,
		},
		{
			"modem not found",
			errors.New("error: couldn't find modem"),
			http.StatusInternalServerError,
			http.StatusNotFound,
			`{"error":"error: couldn't find modem"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buff bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&buff, nil))

			req := newWritersRequest(logger)
			rec := httptest.NewRecorder()

			WriteError(rec, req, test.statusCode, test.err)

			resp := rec.Result()
			body, _ := io.ReadAll(resp.Body)

			if test.expectedStatusCode == http.StatusInternalServerError {
				assert.Contains(t, buff.String(), "level=ERROR msg=Error error=\"something weird happened\"")
			} else {
				if test.err == nil {
					assert.Contains(t, buff.String(), "level=WARN msg=Error error=\"unspecified error\"")
				} else {
					assert.Contains(t, buff.String(), "level=WARN msg=Error error=\""+test.err.Error()+"\"")
				}
			}

			assert.Equal(t, test.expectedStatusCode, resp.StatusCode)
			assert.Equal(t, test.expectedBody, strings.Trim(string(body), " \n"))
		})
	}
}

func newWritersRequest(logger *slog.Logger) *http.Request {
	r := httptest.NewRequest("GET", "/", nil)

	ctx := context.WithValue(r.Context(), "logger", logger)

	return r.WithContext(ctx)
}
