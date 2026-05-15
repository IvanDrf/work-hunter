from datetime import datetime, timezone

from src.core.exc import AccessError, ArgumentError
from src.domain.models.vacancy import VacancyStatus
from src.domain.rules.user import is_user_admin, is_user_employer
from src.domain.rules.vacancy import has_right_to_vacancy, is_vacancy_id_valid
from src.domain.schemas import UserInfo, VacancyCreateSchema, VacancyResponseSchema
from src.infrastructure.service.dependencies import IVacancyRepo
from src.infrastructure.service.dto.vacancy import (
    create_vacancy_dto,
    vacancy_orm_to_response_dto,
)


class VacancyService:
    def __init__(self, vacancy_repo: IVacancyRepo) -> None:
        self.vacancy_repo: IVacancyRepo = vacancy_repo

    async def create_vacancy(self, vacancy: VacancyCreateSchema, user_info: UserInfo) -> VacancyResponseSchema:
        if not user_info.verificated:
            raise AccessError("user is not verificated, can't create vacancy")

        if not is_user_employer(user_info):
            raise AccessError("only employer can create vacancies")

        vacancy_create_date = datetime.now(timezone.utc)
        vacancy_status = VacancyStatus.MODERATING
        vacancyORM = create_vacancy_dto(vacancy, user_info, vacancy_create_date, vacancy_status)

        await self.vacancy_repo.create_vacancy(vacancyORM)

        return vacancy_orm_to_response_dto(vacancyORM)

    async def find_vacancy_by_id(self, vacancy_id: int, user_info: UserInfo | None) -> VacancyResponseSchema | None:
        if not is_vacancy_id_valid(vacancy_id):
            raise ArgumentError(f"vacancy_id must be non negative number, {vacancy_id=}")

        vacancy = await self.vacancy_repo.find_vacancy_by_id(vacancy_id)
        if vacancy is None:
            return None

        if not has_right_to_vacancy(vacancy, user_info):
            raise AccessError("this vacancy is moderating now, you can't see it now")

        return vacancy_orm_to_response_dto(vacancy)

    async def find_vacancies_with_tags(
        self, tags: list[str], offset: int, limit: int, user_info: UserInfo | None
    ) -> list[VacancyResponseSchema] | None:
        if user_info is not None and is_user_admin(user_info):
            vacancies = await self.vacancy_repo.find_vacancies_for_admin_with_tags(tags, offset, limit)
        else:
            vacancies = await self.vacancy_repo.find_only_published_vacancies_with_tags(tags, offset, limit)

        if vacancies is None:
            return None

        return [vacancy_orm_to_response_dto(vacancy) for vacancy in vacancies]

    async def set_vacancy_status(
        self, vacancy_id: int, status: VacancyStatus, moderator_comments: str, user_info: UserInfo
    ) -> None:
        if not is_user_admin(user_info):
            raise AccessError("you can't change vacancy status, you are not admin")

        if not is_vacancy_id_valid(vacancy_id):
            raise ArgumentError(f"vacancy_id must be non negative number, {vacancy_id=}")

        await self.vacancy_repo.set_vacancy_status(vacancy_id, VacancyStatus(status), moderator_comments)

    async def delete_vacancy(self, vacancy_id: int, user_info: UserInfo) -> None:
        if is_user_admin(user_info):
            await self.vacancy_repo.delete_vacancy(vacancy_id)
            return

        author_id = await self.vacancy_repo.find_vacancy_author(vacancy_id)
        if author_id is None:
            raise ArgumentError(f"can't find author for vacancy with {vacancy_id=}")

        if str(author_id) != user_info.user_id:
            raise AccessError("you have no rights to delete vacancy, you didn't created")

        await self.vacancy_repo.delete_vacancy(vacancy_id)
