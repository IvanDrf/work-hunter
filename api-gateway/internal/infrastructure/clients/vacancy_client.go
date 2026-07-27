package clients

import (
	"fmt"
	"log/slog"
	"time"

	"context"

	"github.com/IvanDrf/work-hunter/api-gateway/internal/domain/models"
	"github.com/IvanDrf/work-hunter/api-gateway/internal/infrastructure/adapters"
	"github.com/IvanDrf/work-hunter/pkg/common"
	"github.com/IvanDrf/work-hunter/pkg/vacancy_api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type vacancyClient struct {
	retries int

	conn   *grpc.ClientConn
	client vacancy_api.VacancyClient
}

func NewVacancyClient(host string, port int, retries int) *vacancyClient {
	client, conn := connectToVacancy(host, port)
	return &vacancyClient{
		retries: retries,
		conn:    conn,
		client:  client,
	}
}

func connectToVacancy(host string, port int) (vacancy_api.VacancyClient, *grpc.ClientConn) {
	log := slog.With(slog.String("client", "vacancy"))
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("can't connect to vacancy service", slog.String("error", err.Error()))
		return nil, nil
	}

	return vacancy_api.NewVacancyClient(conn), conn
}

func (c *vacancyClient) Close() {
	log := slog.With(slog.String("client", "vacancy"))

	if c != nil {
		c.conn.Close()
	}

	log.Info("vacancy client is closed")
}

func (c *vacancyClient) Health(ctx context.Context) {
	log := slog.With(slog.String("client", "vacancy"))
	ctx = adapters.InsertLogger(ctx, log)

	resp, err := retry(ctx, c.retries, func() (any, error) {
		resp, err := c.client.Health(ctx, nil)
		if err != nil {
			log.ErrorContext(ctx, "can't check vacancy service health, vacancy service returned error", slog.String("error", err.Error()))
			return nil, err
		}

		return resp, nil
	})

	if err != nil {
		log.ErrorContext(ctx, "vacancy service is not available now", slog.String("error", err.Error()))
	} else {
		log.InfoContext(ctx, "vacancy service is available now", slog.Any("resp", resp))
	}
}

func (c *vacancyClient) CreateVacancy(ctx context.Context, vacancy *models.Vacancy, userInfo *models.UserInfo, companyName string) (*models.VacancyInfo, error) {
	log := slog.With(slog.String("client", "vacancy"))
	ctx = adapters.InsertLogger(ctx, log)

	resp, err := retry(ctx, c.retries, func() (any, error) {
		resp, err := c.client.CreateVacancy(ctx, &vacancy_api.CreateVacancyRequest{
			Title:         vacancy.Title,
			Description:   &vacancy.Description,
			Requirements:  vacancy.Requirements,
			Conditions:    vacancy.Conditions,
			SalaryMin:     vacancy.SalaryMin,
			SalaryMax:     vacancy.SalaryMax,
			Currency:      vacancy_api.Currency(vacancy.Currency),
			City:          vacancy.City,
			Metro:         vacancy.Metro,
			RemoteType:    vacancy_api.RemoteType(vacancy.RemoteType),
			TimeType:      vacancy_api.TimeType(vacancy.TimeType),
			ExperienceMin: vacancy.ExperienceMin,
			ExperienceMax: vacancy.ExperienceMax,
			Tags:          vacancy.Tags,
			UserInfo: &common.FullUserInfo{
				Role:        common.UserRole(userInfo.Role),
				UserId:      userInfo.UserID,
				Verificated: userInfo.Verificated,
				CompanyName: companyName,
			},
		})
		if err != nil {
			log.ErrorContext(ctx, "can't create new vacancy, vacancy service returned error", slog.String("error", err.Error()))
			return nil, err
		}

		return vacancyInfoDTO(resp), nil
	})

	if err != nil {
		return nil, err
	}

	return resp.(*models.VacancyInfo), nil
}

func vacancyInfoDTO(resp *vacancy_api.VacancyInfo) *models.VacancyInfo {
	return &models.VacancyInfo{
		Vacancy: models.Vacancy{
			Title:         resp.GetTitle(),
			Description:   resp.GetDescription(),
			Requirements:  resp.GetRequirements(),
			Conditions:    resp.GetConditions(),
			SalaryMin:     resp.SalaryMin,
			SalaryMax:     resp.SalaryMax,
			Currency:      models.Currency(resp.GetCurrency()),
			City:          resp.City,
			Metro:         resp.Metro,
			RemoteType:    models.RemoteType(resp.GetRemoteType()),
			TimeType:      models.TimeType(resp.GetTimeType()),
			ExperienceMin: resp.ExperienceMin,
			ExperienceMax: resp.ExperienceMax,
			Tags:          resp.GetTags()},

		VacancyID:         resp.GetVacancyId(),
		IsCityValid:       resp.GetIsCityValid(),
		IsMetroValid:      resp.GetIsMetroValid(),
		CreatedAt:         timeDTO(resp.GetCreatedAt()),
		UpdatedAt:         timeDTO(resp.GetUpdatedAt()),
		ClosedAt:          timeDTO(resp.GetClosedAt()),
		Status:            models.VacancyStatus(resp.GetStatus()),
		ModeratedAt:       timeDTO(resp.GetModeratedTime()),
		ModeratorComments: resp.ModeratorComments,
		Views:             resp.Views,
		Applications:      resp.ApplicationsCount,
		AuthorName:        resp.GetAuthorName(),
	}
}

func timeDTO(t *timestamppb.Timestamp) *time.Time {
	if t == nil {
		return nil
	}

	w := t.AsTime()
	return &w
}
