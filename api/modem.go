package api

import (
	"github.com/go-chi/chi/v5"
	utils "go-mmcli-svr/http"
	"net/http"
)

// ModemList
// @Summary	List modems
// @Accept json
// @Produce	json
// @Tags         modem
// @Success	200
// @Router		/modem/ [get]
func (a *Api) ModemList(w http.ResponseWriter, r *http.Request) {
	jsonBuf, err := a.backend.Exec("-L")
	if err != nil {
		utils.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, jsonBuf)
}

// ModemDetail
// @Summary	Modem info
// @Accept json
// @Produce	json
// @Tags         modem
// @Success	200
// @Failure	404 {object} http.ErrorResponse
// @Router		/modem/{id} [get]
// @Param	id		path	string  true  "Modem-Id"
func (a *Api) ModemDetail(w http.ResponseWriter, r *http.Request) {
	id, err := ValidatePathIndex(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, r, http.StatusBadRequest, err)
		return
	}

	jsonBuf, err := a.backend.Exec("-m", id)
	if err != nil {
		utils.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, jsonBuf)
}
