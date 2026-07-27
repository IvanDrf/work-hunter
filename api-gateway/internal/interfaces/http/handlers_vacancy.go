package http

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/IvanDrf/work-hunter/api-gateway/internal/domain/models"
	"github.com/IvanDrf/work-hunter/api-gateway/internal/infrastructure/adapters"
)

func (h *Handlers) CreateVacancy(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.requestTime)
	defer cancel()

	log := slog.With(slog.String("handlers", "CreateVacancy"))
	log.InfoContext(ctx, "request")

	ctx = adapters.InsertLogger(ctx, log)

	userInfo, err := getUserInfo(ctx)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	vacancy := &models.Vacancy{}
	if err = json.NewDecoder(r.Body).Decode(vacancy); err != nil {
		log.InfoContext(ctx, "can't parse requests's body", slog.String("error", err.Error()))
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(models.Error{
			Message: invalidBodyRequestMessage,
			Code:    models.ErrCodeUnprocessableEntity,
		})
		return
	}
	defer r.Body.Close()

	companyName := ""
	// panic("send request to user service for company name")

	resp, err := h.vacancyClient.CreateVacancy(ctx, vacancy, userInfo, companyName)
	if err != nil {
		handleResponseError(w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
