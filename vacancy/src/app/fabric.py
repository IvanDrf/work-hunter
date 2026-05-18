from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from src.api.handlers import VacancyHandlers
from src.core.config.config import Config
from src.infrastructure.persistence.postgresql_repo import TagRepo, UnitOfWork, VacancyRepo
from src.infrastructure.persistence.postgresql_repo import connect as connect_postgresql
from src.infrastructure.service.dependencies import ITagRepo, IUnitOfWork, IVacancyRepo
from src.infrastructure.service.vacancy import VacancyService


class Fabric:
    def __init__(self, config: Config) -> None:
        self.config: Config = config

    async def new_handlers(self) -> VacancyHandlers:
        session_maker = await connect_postgresql(self.config.database)

        vacancy_repo: IVacancyRepo = self.new_vacancy_repo()
        tag_repo: ITagRepo = self.new_tag_repo()
        uof: IUnitOfWork = self.new_unit_of_work(session_maker)

        vacancy_service = self.new_vacancy_service(vacancy_repo, tag_repo, uof)

        return VacancyHandlers(vacancy_service)

    def new_vacancy_repo(self) -> VacancyRepo:
        return VacancyRepo()

    def new_tag_repo(self) -> TagRepo:
        return TagRepo()

    def new_unit_of_work(self, session_maker: async_sessionmaker[AsyncSession]) -> UnitOfWork:
        return UnitOfWork(session_maker)

    def new_vacancy_service(self, vacancy_repo: IVacancyRepo, tag_repo: ITagRepo, uof: IUnitOfWork) -> VacancyService:
        return VacancyService(vacancy_repo, tag_repo, uof)
