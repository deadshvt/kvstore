package user

import (
	"encoding/json"
	"net/http"

	"github.com/deadshvt/kvstore/internal/entity"
	"github.com/deadshvt/kvstore/internal/errs"
	"github.com/deadshvt/kvstore/pkg/response"
)

const (
	WrongCredentials = "Wrong credentials"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info().Msg("Logging in...")

	var req LoginRequestDTO
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		msg := errs.WrapError(errs.ErrJSONDecode, err).Error()
		h.Logger.Error().Msg(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	user := &entity.User{
		Username: req.Username,
		Password: req.Password,
	}

	result, err := h.Usecase.Login(r.Context(), user)
	if err != nil {
		msg := errs.WrapError(errs.ErrLogin, err).Error()
		h.Logger.Error().Msg(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	resp := &LoginResponseDTO{}

	statusCode := http.StatusOK
	if !result.Success {
		resp.Error = WrongCredentials
		statusCode = http.StatusUnauthorized
	} else {
		resp.Token = result.Token
	}

	response.Set(w, statusCode, resp)
	h.Logger.Info().Msg("Logged in")
}

type LoginRequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponseDTO struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}
