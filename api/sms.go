package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	utils "go-mmcli-svr/http"
	"net/http"
	"strings"
)

type SmsRequestData struct {
	Number string `json:"number" example:"+431234567890"`
	Text   string `json:"text" example:"Ping"`
}

// SmsCreated {"modem":{"messaging":{"created-sms":"/org/freedesktop/ModemManager1/SMS/0"}}}
type SmsCreated struct {
	Modem struct {
		Messaging struct {
			CreatedSms string `json:"created-sms"`
		} `json:"messaging"`
	} `json:"modem"`
}

// SmsGet
// @Summary	List SMS messages
// @Accept json
// @Produce	json
// @Tags         sms
// @Success	200
// @Failure	404 {object} http.ErrorResponse
// @Router		/sms/ [get]
// @Router		/sms/{id} [get]
// @Param	id		path	string  true  "SMS-Id"
// @Param	modem	query	string  false  "Modem-Id"
func (a *Api) SmsGet(w http.ResponseWriter, r *http.Request) {
	modem := r.Context().Value("modem").(string)

	id, err := ValidatePathIndex(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, r, http.StatusBadRequest, err)
		return
	}

	var jsonBuf []byte

	if len(id) > 0 {
		jsonBuf, err = a.backend.Exec("-s", id)
	} else {
		jsonBuf, err = a.backend.ExecModem(modem, "--messaging-list-sms")
	}

	if err != nil {
		utils.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, jsonBuf)
}

// SmsSend
// @Summary	Send SMS
// @Accept json
// @Produce	json
// @Tags         sms
// @Success	200
// @Param        data   body      SmsRequestData  true  "Data"
// @Router		/sms [post]
// @Param	modem	query	string  false  "Modem-Id"
func (a *Api) SmsSend(w http.ResponseWriter, r *http.Request) {
	modem := r.Context().Value("modem").(string)

	var requestData SmsRequestData
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		utils.WriteError(w, r, http.StatusBadRequest, err)
		return
	}

	smsStr := fmt.Sprintf("number='%s',text='%s'", requestData.Number, requestData.Text)
	jsonBuf, err := a.backend.ExecModem(modem, "--messaging-create-sms="+smsStr)
	if err != nil {
		utils.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	var smsCreated SmsCreated
	err = json.NewDecoder(bytes.NewReader(jsonBuf)).Decode(&smsCreated)
	if err != nil {
		utils.WriteError(w, r, http.StatusBadRequest, err)
		return
	}

	jsonBuf, err = a.backend.Exec("-s", smsCreated.Modem.Messaging.CreatedSms, "--send")
	if err != nil {
		utils.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	resp := Response{
		Message: strings.Trim(string(jsonBuf), "\n "),
	}
	utils.WriteJson(w, http.StatusOK, resp)
}

// SmsDelete
// @Summary	Delete a single SMS
// @Accept json
// @Produce	json
// @Tags         sms
// @Success	200
// @Router		/sms/{id} [delete]
// @Param	id		path	string  true  "SMS-Id"
func (a *Api) SmsDelete(w http.ResponseWriter, r *http.Request) {
	modem := r.Context().Value("modem").(string)

	id, err := ValidatePathIndex(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, r, http.StatusBadRequest, err)
		return
	}

	jsonBuf, err := a.backend.ExecModem(modem, "--messaging-delete-sms="+id)
	if err != nil {
		utils.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	resp := Response{
		Message: strings.Trim(string(jsonBuf), "\n "),
	}
	utils.WriteJson(w, http.StatusOK, resp)
}
