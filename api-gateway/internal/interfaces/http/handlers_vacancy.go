package http

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

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

	companyName, err := h.userClient.SendGetCompanyName(ctx, userInfo.UserID)
	if err != nil {
		handleResponseError(w, err)
		return
	}

	resp, err := h.vacancyClient.SendCreateVacancy(ctx, vacancy, userInfo, companyName)
	if err != nil {
		handleResponseError(w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handlers) FindVacancyByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.requestTime)
	defer cancel()

	log := slog.With(slog.String("handlers", "FindVacancyByID"))
	log.InfoContext(ctx, "request")

	ctx = adapters.InsertLogger(ctx, log)

	id := r.URL.Query().Get("vacancy")
	if id == "" {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(models.Error{
			Message: "vacancy id is empty, invalid url, needs query arg",
			Code:    models.ErrCodeInvalidArgument,
		})
		return
	}

	vacancyID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		log.Info("invalid vacancy id in query", slog.String("error", err.Error()))
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Error{
			Message: "invalid vacancy id in query, not uint64",
			Code:    models.ErrCodeInvalidArgument,
		})
		return
	}

	userInfo, _ := getUserInfo(ctx) // ignore error, because user info can be nil and in this request it's valid

	resp, err := h.vacancyClient.SendFindVacancyByID(ctx, vacancyID, userInfo)
	if err != nil {
		handleResponseError(w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
