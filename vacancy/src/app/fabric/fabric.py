from src.api.handlers import VacancyHandlers
from src.app.fabric.vacancy_repo import new_vacancy_repo
from src.app.fabric.vacancy_service import new_vacancy_service
from src.core.config.config import Config


class Fabric:
    def __init__(self, config: Config) -> None:
        self.config: Config = config

    def new_handlers(self) -> VacancyHandlers:
        vacancy_repo = new_vacancy_repo(self.config.database)
        vacancy_service = new_vacancy_service(vacancy_repo)

        return VacancyHandlers(vacancy_service)
