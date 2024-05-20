package api

import (
	"github.com/mbretter/go-mmcli-svr/backend"
	"testing"
)

func ProvideTestApi(t *testing.T) (*Api, *backend.BackendMock) {
	backendMock := backend.NewBackendMock(t)

	return Provide(backendMock), backendMock
}
