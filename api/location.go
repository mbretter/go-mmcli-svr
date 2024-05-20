package api

import (
	utils "github.com/mbretter/go-mmcli-svr/http"
	"net/http"
)

// LocationGet
// @Summary	Get location
// @Accept json
// @Produce	json
// @Tags location
// @Success	200
// @Failure	404 {object} http.ErrorResponse
// @Router		/location [get]
// @Param	modem		query	string  false  "Modem-Id"
func (a *Api) LocationGet(w http.ResponseWriter, r *http.Request) {
	modem := r.Context().Value("modem").(string)

	jsonBuf, err := a.backend.ExecModem(modem, "--location-get")
	if err != nil {
		utils.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, jsonBuf)
}
