from src.infrastructure.persistence.postgresql_repo.vacancy_dml import VacancyDMLRepo
from src.infrastructure.persistence.postgresql_repo.vacancy_search import VacancySearchRepo


class VacancyRepo(VacancySearchRepo, VacancyDMLRepo):
    pass
