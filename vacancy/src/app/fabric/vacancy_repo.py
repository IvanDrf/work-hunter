from src.core.config.database import PostgreSQLConfig
from src.database.postgresql import connect
from src.repo.vacancy import VacancyRepo


def new_vacancy_repo(config: PostgreSQLConfig) -> VacancyRepo:
    return VacancyRepo(connect(config))
