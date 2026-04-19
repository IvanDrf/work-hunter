package handlers

import (
	"context"
	"errors"
	"log/slog"

	user_api "github.com/IvanDrf/work-hunter/pkg/user-api"
	"github.com/IvanDrf/work-hunter/users/internal/domain/models"
	"github.com/IvanDrf/work-hunter/users/internal/interfaces/grpc/dto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) ListUsers(ctx context.Context, req *user_api.ListUsersRequest) (*user_api.ListUsersResponse, error) {
	log := h.log.With(slog.String("scope", "intarfaces/grpc/handlers/ListUsers"))
	log.Info("ListUsers got request")

	resp, err := h.UserService.ListUsers(ctx, &dto.ListUsersRequest{
		PageSize:    req.PageSize,
		Status:      user_api.UserStatus_name[int32(req.Status)],
		Role:        user_api.UserRole_name[int32(req.Role)],
		SearchQuery: req.SerchQuery,
		SortBy:      req.SortBy,
		// TODO: regenerate pb files
		Offset: 0,
	})

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

	log.Info("List users successfully response")
	return convertListDtoToListResp(resp), nil
}
