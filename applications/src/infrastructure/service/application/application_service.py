import logging
from asyncio import gather, wait_for
from uuid import uuid4

from src.core.exc import AccessError, AlreadyExistsError, NotFoundError
from src.domain.schemas import ApplicationMessage, ApplicationSchema, UserInfo, UserRole
from src.infrastructure.service.application.dependencies import IApplicationProducer, IApplicationRepo, IUnitOfWork
from src.infrastructure.service.application.dto import application_dto


class ApplicationService:
    def __init__(
        self,
        uof: IUnitOfWork,
        application_repo: IApplicationRepo,
        repo_timeout: float,
        application_producer: IApplicationProducer,
    ) -> None:
        self.uof: IUnitOfWork = uof
        self.application_repo: IApplicationRepo = application_repo
        self.repo_timeout: float = repo_timeout

        self.application_producer: IApplicationProducer = application_producer

        self.logger = logging.getLogger("ApplicationService")

    async def stop(self) -> None:
        await gather(*[self.uof.stop(), self.application_producer.stop()])

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
            await self.application_producer.publish_application(
                ApplicationMessage(message_id=uuid4(), vacancy_id=application.vacancy_id)
            )

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


def count_applications_amount(messages: list[ApplicationSchema | None]) -> dict[int, int]:
    amounts = {}
    for message in messages:
        if message is not None:
            amounts[message.vacancy_id] = amounts.get(message.vacancy_id, 0) + 1

    return amounts
