from datetime import timedelta

from src.infrastructure.service.dependencies import ICache, ITagRepo, IUnitOfWork, IVacancyRepo, IValidationServiceClient
from src.infrastructure.service.vacancy_dml import VacancyDMLService
from src.infrastructure.service.vacancy_search import VacancySearchService


class VacancyService(VacancyDMLService, VacancySearchService):
    def __init__(
        self,
        vacancy_repo: IVacancyRepo,
        tag_repo: ITagRepo,
        uof: IUnitOfWork,
        cache: ICache,
        validation_client: IValidationServiceClient,
        vacancy_ttl: timedelta,
        cache_timeout: float,
    ) -> None:
        VacancyDMLService.__init__(self, vacancy_repo, tag_repo, uof, cache, validation_client, vacancy_ttl, cache_timeout)
        VacancySearchService.__init__(self, vacancy_repo, uof, cache, vacancy_ttl, cache_timeout)

    async def stop(self) -> None:
        await self.uof.stop()
