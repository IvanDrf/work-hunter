from asyncio import wait_for

from src.core.exc import AccessError, AlreadyExistsError
from src.domain.schemas import ApplicationSchema, UserRole
from src.infrastructure.service.application.dependencies import IApplicationRepo, IUnitOfWork
from src.infrastructure.service.application.dto import application_dto


class ApplicationService:
    def __init__(self, uof: IUnitOfWork, application_repo: IApplicationRepo, repo_timeout: float) -> None:
        self.uof: IUnitOfWork = uof
        self.application_repo: IApplicationRepo = application_repo
        self.repo_timeout: float = repo_timeout

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
