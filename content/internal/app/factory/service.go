package factory

import (
	"github.com/IvanDrf/work-hunter/content/internal/domain/ports/repo"
	"github.com/IvanDrf/work-hunter/content/internal/domain/ports/service"
	"github.com/IvanDrf/work-hunter/content/internal/infrastructure/services"
)

func (f *Factory) newContentService(repo repo.ContentRepo) service.ContentService {
	return services.NewContentService(repo, f.log)
}
