package pair

import (
	"encoding/json"
	"net/http"

	"github.com/deadshvt/kvstore/internal/entity"
	"github.com/deadshvt/kvstore/internal/errs"
	"github.com/deadshvt/kvstore/pkg/response"
)

const (
	StatusSuccess = "success"
	StatusFail    = "fail"
)

func (h *Handler) SetPairs(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info().Msg("Setting pairs...")

	var req SetPairsRequestDTO
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		msg := errs.WrapError(errs.ErrJSONDecode, err).Error()
		h.Logger.Error().Msg(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if len(req.Data) == 0 {
		msg := errs.WrapError(errs.ErrEmptyPairs, err).Error()
		h.Logger.Error().Msg(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	pairs := make([]*entity.Pair, 0, len(req.Data))
	for key, value := range req.Data {
		pairs = append(pairs, &entity.Pair{
			Key:   key,
			Value: value,
		})
	}

	result, err := h.Usecase.SetPairs(r.Context(), pairs)
	if err != nil {
		msg := errs.WrapError(errs.ErrSetPairs, err).Error()
		h.Logger.Error().Msg(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	resp := &SetPairsResponseDTO{
		Status: StatusSuccess,
	}

	statusCode := http.StatusOK

	if !result.Success {
		resp.Status = StatusFail
		resp.Errors = make(map[string]string, len(result.Errors))
		for _, err := range result.Errors {
			resp.Errors[err.Key] = err.Message
		}

		statusCode = http.StatusInternalServerError
	}

	response.Set(w, statusCode, resp)
	h.Logger.Info().Msg("Set pairs")
}

type SetPairsRequestDTO struct {
	Data map[string]interface{} `json:"data"`
}

type SetPairsResponseDTO struct {
	Status string            `json:"status"`
	Errors map[string]string `json:"errors,omitempty"`
}
