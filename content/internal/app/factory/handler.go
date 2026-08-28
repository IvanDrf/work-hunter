package factory

import (
	"github.com/IvanDrf/work-hunter/content/internal/interfaces/http/handlers"
)

func (f *Factory) NewHandler() (*handlers.ContentHandler, error) {
	repo, err := f.newRepos()
	if err != nil {
		return nil, err
	}

	service := f.newContentService(repo)

	return handlers.NewContentHandler(service), nil
}
