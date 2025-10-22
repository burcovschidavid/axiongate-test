package handlers

import (
	"encoding/json"
	"net/http"
	"shipping-api/internal/core/domain"
	"shipping-api/internal/core/ports"
)

type ShippingHandler struct {
	service ports.ShippingService
}

func NewShippingHandler(service ports.ShippingService) *ShippingHandler {
	return &ShippingHandler{
		service: service,
	}
}

func (h *ShippingHandler) CreateShipment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request domain.GenericShippingRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	provider := r.URL.Query().Get("provider")

	if provider == "" {
		responses, err := h.service.BroadcastShipment(r.Context(), &request)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, responses)
		return
	}

	response, err := h.service.ProcessShipment(r.Context(), &request, provider)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
