package api

import (
	"net/http"
)

type Utils struct {
}

func ProvideUtilsApi() *Utils {
	return &Utils{}
}

// Ping
// @Summary	ping
// @Accept plain
// @Produce	plain
// @Tags         utils
// @Success	200	{string} boing
// @Router		/ping [get]
// @Security JWT
func (u *Utils) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	//goland:noinspection ALL
	w.Write([]byte("boing"))
}
