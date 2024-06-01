package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/mbretter/go-mmcli-svr/backend"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterSmsRoutes(t *testing.T) {
	backendMock := backend.NewBackendMock(t)
	r := chi.NewRouter()
	registerSmsRoutes(r, backendMock)

	routes := r.Routes()
	assert.Len(t, routes, 3)
	for _, route := range routes {
		switch route.Pattern {
		case "/sms/":
			assert.Len(t, route.Handlers, 1)
		case "/sms/{id:[a-zA-Z0-9%/]+}":
			assert.Len(t, route.Handlers, 2)
		case "/sms":
			assert.Len(t, route.Handlers, 1)
		default:
			t.Errorf("Unexpected route: %s", route.Pattern)
		}
	}
}
