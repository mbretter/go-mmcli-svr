package api

import (
	"go-mmcli-svr/backend"
	"testing"
)

func ProvideTestApi(t *testing.T) (*Api, *backend.BackendMock) {
	backendMock := backend.NewBackendMock(t)

	return Provide(backendMock), backendMock
}
