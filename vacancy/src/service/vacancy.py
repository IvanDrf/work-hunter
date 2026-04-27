from pkg.common.common_pb2 import UserInfo, UserRole
from pkg.vacancy_api.vacancy_pb2 import VacancyInfo

from src.core.exc import AccessError, ArgumentError
from src.domain.models.vacancy import VacancyORM, VacancyStatus
from src.domain.rules.vacancy import check_vacancy_fields
from src.service.dependencies.repo import IVacancyRepo
from src.service.dto.vacancy import create_vacancy_dto, vacancy_info_dto


class VacancyService:
    def __init__(self, vacancy_repo: IVacancyRepo) -> None:
        self.vacancy_repo: IVacancyRepo = vacancy_repo

    async def create_vacancy(self, vacancy: VacancyInfo, user_info: UserInfo) -> int:
        '''
        Raises:
            AccessError: if user is not verificated
            ArgumentError: if field of vacancy is invalid
            InternalError: from vacancy_repo
        '''

        if not user_info.verificated:
            raise AccessError(
                '''user is not verificated, can't create vacancy '''
            )

        check_vacancy_fields(vacancy)

        vacancy_id = await self.vacancy_repo.create_vacancy(create_vacancy_dto(vacancy, user_info))
        return vacancy_id

    async def find_vacancy_by_id(self, vacancy_id: int, user_info: UserInfo) -> VacancyInfo | None:
        '''
        Raises:
            ArgumentError: if vacancy_id < 0
            AccessError: if vacancy is moderating now and user trying to find vacancy but he is not admin
            InternalError: from vacancy_repo
        '''

        if vacancy_id < 0:
            raise ArgumentError(
                f'vacancy must be non negative number, {vacancy_id=}'
            )

        vacancy = await self.vacancy_repo.find_vacancy_by_id(vacancy_id)
        if vacancy is None:
            return None

        if not has_right_to_vacancy(vacancy, user_info):
            raise AccessError(
                '''this vacancy is moderating now, you can't see it now '''
            )

        return vacancy_info_dto(vacancy)


def has_right_to_vacancy(vacancy: VacancyORM, user_info: UserInfo) -> bool:
    return vacancy.status == VacancyStatus.MODERATING and user_info.role != UserRole.ADMIN and vacancy.author_id != user_info.user_id
