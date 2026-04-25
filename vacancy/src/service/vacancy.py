from pkg.common.common_pb2 import UserInfo
from pkg.vacancy_api.vacancy_pb2 import VacancyInfo

from src.core.exc.external import ExternalError
from src.domain.rules.vacancy import check_vacancy_fields
from src.service.dependencies.repo import IVacancyRepo
from src.service.dto.vacancy import create_vacancy_dto, vacancy_info_dto


class VacancyService:
    def __init__(self, vacancy_repo: IVacancyRepo) -> None:
        self.vacancy_repo: IVacancyRepo = vacancy_repo

    async def create_vacancy(self, vacancy: VacancyInfo, user_info: UserInfo) -> int:
        '''
        Raises:
            ExternalError: from vacancy field validation
            InternalError: from vacancy_repo
        '''
        if not user_info.verificated:
            raise ExternalError(
                '''user is not verificated, can't create vacancy '''
            )

        check_vacancy_fields(vacancy)

        vacancy_id = await self.vacancy_repo.create_vacancy(create_vacancy_dto(vacancy, user_info))
        return vacancy_id

    async def find_vacancy_by_id(self, vacancy_id: int) -> VacancyInfo | None:
        '''
        Raises:
            ExternalError: if vacancy_id < 0
            InternalError: from vacancy_repo
        '''

        if vacancy_id < 0:
            raise ExternalError(
                f'vacancy must be non negative number, {vacancy_id=}'
            )

        vacancy = await self.vacancy_repo.find_vacancy_by_id(vacancy_id)
        return None if vacancy is None else vacancy_info_dto(vacancy)
