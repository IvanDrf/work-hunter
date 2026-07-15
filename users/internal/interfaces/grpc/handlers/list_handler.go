package handlers

import (
	"context"
	"errors"
	"log/slog"

	user_api "github.com/IvanDrf/work-hunter/pkg/user-api"
	"github.com/IvanDrf/work-hunter/users/internal/domain/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) ListUsers(ctx context.Context, req *user_api.ListUsersRequest) (*user_api.ListUsersResponse, error) {
	log := h.log.With(slog.String("scope", "handler/ListUsers"))
	log.Info("ListUsers called")

	resp, err := h.UserService.ListUsers(ctx, convertListReqToDto(req))

	if err != nil {
		var e models.Error
		if errors.As(err, &e) {
			switch e.Code {
			case models.ErrCodeInternal:
				return nil, status.Error(codes.Internal, e.Message)
			case models.ErrCodeUserNotFound:
				return nil, status.Error(codes.NotFound, e.Message)
			default:
				return nil, status.Error(codes.InvalidArgument, e.Message)
			}
		}
		log.Error("unhandled error", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "internal server error")
	}

	// Конвертация ответа
	return convertListDtoToListResp(resp), nil
}
