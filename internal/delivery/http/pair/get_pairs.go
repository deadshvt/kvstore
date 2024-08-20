package pair

import (
	"encoding/json"
	"net/http"

	"github.com/deadshvt/kvstore/internal/errs"
	"github.com/deadshvt/kvstore/pkg/response"
)

func (h *Handler) GetPairs(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info().Msg("Getting pairs...")

	var req GetPairsRequestDTO
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		msg := errs.WrapError(errs.ErrJSONDecode, err).Error()
		h.Logger.Error().Msg(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	result, err := h.Usecase.GetPairs(r.Context(), req.Keys)
	if err != nil {
		msg := errs.WrapError(errs.ErrGetPairs, err).Error()
		h.Logger.Error().Msg(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	data := make(map[string]interface{}, len(result.Pairs))
	for _, pair := range result.Pairs {
		data[pair.Key] = pair.Value
	}

	errors := make(map[string]string, len(result.Errors))
	for _, err := range result.Errors {
		errors[err.Key] = err.Message
	}

	resp := &GetPairsResponseDTO{
		Data:   data,
		Errors: errors,
	}

	statusCode := http.StatusOK
	if len(result.Errors) > 0 {
		statusCode = http.StatusInternalServerError
	}

	response.Set(w, statusCode, resp)
	h.Logger.Info().Msg("Got pairs")
}

type GetPairsRequestDTO struct {
	Keys []string `json:"keys"`
}

type GetPairsResponseDTO struct {
	Data   map[string]interface{} `json:"data"`
	Errors map[string]string      `json:"errors,omitempty"`
}
