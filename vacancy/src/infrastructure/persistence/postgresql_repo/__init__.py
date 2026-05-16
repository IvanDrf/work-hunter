from src.infrastructure.persistence.postgresql_repo.connection import connect
from src.infrastructure.persistence.postgresql_repo.vacancy import VacancyRepo

__all__ = [
    "connect",
    "VacancyRepo",
]
