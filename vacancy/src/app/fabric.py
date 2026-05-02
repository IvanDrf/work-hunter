from src.api.handlers import VacancyHandlers
from src.core.config.config import Config
from src.database.postgresql import connect
from src.infrastructure.repo.vacancy import VacancyRepo
from src.infrastructure.service.dependencies.repo import IVacancyRepo
from src.infrastructure.service.vacancy import VacancyService


class Fabric:
    def __init__(self, config: Config) -> None:
        self.config: Config = config

    async def new_handlers(self) -> VacancyHandlers:
        vacancy_repo = await self.new_vacancy_repo()
        vacancy_service = self.new_vacancy_service(vacancy_repo)

        return VacancyHandlers(vacancy_service)

    async def new_vacancy_repo(self) -> VacancyRepo:
        return VacancyRepo(await connect(self.config.database))

    def new_vacancy_service(self, repo: IVacancyRepo) -> VacancyService:
        return VacancyService(repo)
