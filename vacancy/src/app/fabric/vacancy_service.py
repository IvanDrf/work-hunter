from src.service.dependencies.repo import IVacancyRepo
from src.service.vacancy import VacancyService


def new_vacancy_service(repo: IVacancyRepo) -> VacancyService:
    return VacancyService(repo)
