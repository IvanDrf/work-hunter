from pkg.common.common_pb2 import UserInfo
from pkg.vacancy_api.vacancy_pb2 import Vacancies, VacancyInfo
from pkg.vacancy_api.vacancy_pb2 import VacancyStatus as PBVacancyStatus

from src.core.exc import AccessError, ArgumentError
from src.domain.models.vacancy import VacancyStatus
from src.domain.rules.user import is_user_admin, is_user_employer
from src.domain.rules.vacancy import check_vacancy_fields, has_right_to_vacancy, is_vacancy_id_valid
from src.service.dependencies.repo import IVacancyRepo
from src.service.dto.vacancy import create_vacancy_dto, find_vacancies_with_tags_dto, vacancy_info_dto


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

        if not is_user_employer(user_info):
            raise AccessError(
                '''only employer can create vacancies'''
            )

        check_vacancy_fields(vacancy)

        vacancy_id = await self.vacancy_repo.create_vacancy(create_vacancy_dto(vacancy, user_info))
        return vacancy_id

    async def find_vacancy_by_id(self, vacancy_id: int, user_info: UserInfo | None) -> VacancyInfo | None:
        '''
        Raises:
            ArgumentError: if vacancy_id < 0
            AccessError: if vacancy is moderating now and user trying to find vacancy but he is not admin
            InternalError: from vacancy_repo
        '''

        if not is_vacancy_id_valid(vacancy_id):
            raise ArgumentError(
                f'vacancy_id must be non negative number, {vacancy_id=}'
            )

        vacancy = await self.vacancy_repo.find_vacancy_by_id(vacancy_id)
        if vacancy is None:
            return None

        if not has_right_to_vacancy(vacancy, user_info):
            raise AccessError('''this vacancy is moderating now, you can't see it now ''')

        return vacancy_info_dto(vacancy)

    async def find_vacancies_with_tags(self, tags: list[str], offset: int, limit: int, user_info: UserInfo | None) -> Vacancies | None:
        '''
        Raises:
            ArgumentError: if offset < MIN_OFFSET, limit is invalid, amount of tags is invalid
            InternalError: from vacancy_repo
        '''

        if user_info is not None and is_user_admin(user_info):
            vacancies = await self.vacancy_repo.find_vacancies_for_admin_with_tags(tags, offset, limit)
        else:
            vacancies = await self.vacancy_repo.find_only_published_vacancies_with_tags(tags, offset, limit)

        if vacancies is None:
            return None

        return Vacancies(vacancies=find_vacancies_with_tags_dto(vacancies), limit=limit, offset=offset)

    async def set_vacancy_status(self, vacancy_id: int, status: PBVacancyStatus, moderator_comments: str, user_info: UserInfo) -> None:
        if not is_user_admin(user_info):
            raise AccessError(
                '''you can't change vacancy status, you are not admin'''
            )

        if not is_vacancy_id_valid(vacancy_id):
            raise ArgumentError(
                f'vacancy_id must be non negative number, {vacancy_id=}'
            )

        await self.vacancy_repo.set_vacancy_status(vacancy_id, VacancyStatus(status), moderator_comments)

    async def delete_vacancy(self, vacancy_id: int, user_info: UserInfo) -> None:
        if is_user_admin(user_info):
            await self.vacancy_repo.delete_vacancy(vacancy_id)
            return

        author_id = await self.vacancy_repo.find_vacancy_author(vacancy_id)
        if author_id is None:
            raise ArgumentError(
                f'''can't find author for vacancy with {vacancy_id=}'''
            )

        if str(author_id) != user_info.user_id:
            raise AccessError(
                f'''you have no rights to delete vacancy, you didn't created'''
            )

        await self.vacancy_repo.delete_vacancy(vacancy_id)
