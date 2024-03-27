package userbalance

import (
	"errors"
	"net/http"

	"github.com/citadel-corp/paimon-bank/internal/common/middleware"
	"github.com/citadel-corp/paimon-bank/internal/common/request"
	"github.com/citadel-corp/paimon-bank/internal/common/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var req CreateUserBalancePayload

	err = request.DecodeJSON(w, r, &req)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Failed to decode JSON",
			Error:   err.Error(),
		})
		return
	}

	req.UserID = userID

	err = req.Validate()
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: err.Error(),
		})
		return
	}

	resp := h.service.Create(r.Context(), req)
	if resp.Error != "" {
		response.JSON(w, resp.Code, response.ResponseBody{
			Message: resp.Message,
			Error:   resp.Error,
		})
		return
	}

	response.JSON(w, resp.Code, response.ResponseBody{
		Message: resp.Message,
	})
}

func (h *Handler) Transaction(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var req CreateTransactionPayload

	err = request.DecodeJSON(w, r, &req)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Failed to decode JSON",
			Error:   err.Error(),
		})
		return
	}

	req.UserID = userID

	err = req.Validate()
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: err.Error(),
		})
		return
	}

	resp := h.service.CreateTransaction(r.Context(), req)
	if resp.Error != "" {
		response.JSON(w, resp.Code, response.ResponseBody{
			Message: resp.Message,
			Error:   resp.Error,
		})
		return
	}

	response.JSON(w, resp.Code, response.ResponseBody{
		Message: resp.Message,
	})
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var req ListUserBalancePayload

	req.UserID = userID

	resp := h.service.List(r.Context(), req)
	if resp.Error != "" {
		response.JSON(w, resp.Code, response.ResponseBody{
			Message: resp.Message,
			Error:   resp.Error,
		})
		return
	}

	response.JSON(w, resp.Code, response.ResponseBody{
		Message: resp.Message,
		Data:    resp.Data,
	})
}

func getUserID(r *http.Request) (string, error) {
	if authValue, ok := r.Context().Value(middleware.ContextAuthKey{}).(string); ok {
		return authValue, nil
	}

	return "", errors.New("unauthorized")
}
