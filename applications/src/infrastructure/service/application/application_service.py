from asyncio import wait_for

from src.core.exc import AccessError, AlreadyExistsError, NotFoundError
from src.domain.schemas import ApplicationSchema, UserInfo, UserRole
from src.infrastructure.service.application.dependencies import IApplicationRepo, IUnitOfWork
from src.infrastructure.service.application.dto import application_dto


class ApplicationService:
    def __init__(self, uof: IUnitOfWork, application_repo: IApplicationRepo, repo_timeout: float) -> None:
        self.uof: IUnitOfWork = uof
        self.application_repo: IApplicationRepo = application_repo
        self.repo_timeout: float = repo_timeout

    async def stop(self) -> None:
        await self.uof.stop()

    async def update_application(self, application: ApplicationSchema) -> None:
        if not application.user_info.verificated:
            raise AccessError("only verificated users can apply for a job")

        if application.user_info.user_role != UserRole.EMPLOYEE:
            raise AccessError("only employee can apply for a job")

        async with self.uof as uof:
            app = await wait_for(
                self.application_repo.find_application(uof, application.vacancy_id, application.user_info.user_id),
                timeout=self.repo_timeout,
            )

            if app is not None:
                raise AlreadyExistsError("you have already applied for this job")

            await self.application_repo.add_application(uof, application_dto(application))

    async def find_vacancies_ids_by_applications(self, user_info: UserInfo, *, limit: int, offset: int) -> list[int]:
        if user_info.user_role != UserRole.EMPLOYEE and user_info.user_role != UserRole.ADMIN:
            raise AccessError("only employee can apply for a job, so there is no vacancies")

        async with self.uof as uof:
            vacancies_ids = await wait_for(
                self.application_repo.find_vacancies_ids_by_user_id(
                    uof,
                    user_info.user_id,
                    limit=limit,
                    offset=offset,
                ),
                timeout=self.repo_timeout,
            )

        if vacancies_ids is None:
            raise NotFoundError(f"can't find any vacancies for user_id={user_info.user_id}, {limit=}, {offset=}")

        return vacancies_ids
