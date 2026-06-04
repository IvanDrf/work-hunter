from datetime import timedelta

from src.infrastructure.service.dependencies.cache import ICache
from src.infrastructure.service.dependencies.repo import ITagRepo, IVacancyRepo
from src.infrastructure.service.dependencies.unit_of_work import IUnitOfWork
from src.infrastructure.service.vacancy_dml import VacancyDMLService
from src.infrastructure.service.vacancy_search import VacancySearchService


class VacancyService(VacancyDMLService, VacancySearchService):
    def __init__(
        self,
        vacancy_repo: IVacancyRepo,
        tag_repo: ITagRepo,
        uof: IUnitOfWork,
        cache: ICache,
        vacancy_ttl: timedelta,
        cache_timeout: float,
    ) -> None:
        VacancyDMLService.__init__(self, vacancy_repo, tag_repo, uof, cache, vacancy_ttl, cache_timeout)
        VacancySearchService.__init__(self, vacancy_repo, uof, cache, vacancy_ttl, cache_timeout)
