import logging

from pkg.common.common_pb2 import UserInfo
from pkg.vacancy_api.vacancy_pb2 import VacancyInfo

from src.core.exc.external import ExternalError
from src.domain.rules.vacancy import check_vacancy_fields
from src.service.dependencies.repo import IVacancyRepo
from src.service.dto.vacancy import create_vacancy_dto


class VacancyService:
    def __init__(self, vacancy_repo: IVacancyRepo) -> None:
        self.vacancy_repo: IVacancyRepo = vacancy_repo

    async def create_vacancy(self, vacancy: VacancyInfo, user_info: UserInfo) -> int:
        try:
            check_vacancy_fields(vacancy)

            vacancy_id = await self.vacancy_repo.create_vacancy(create_vacancy_dto(vacancy, user_info))
            return vacancy_id
        except ExternalError as e:
            logging.info(f'create_vacancy: {e}')

            raise ExternalError(f'Invalid field in vacancy: {e}')
