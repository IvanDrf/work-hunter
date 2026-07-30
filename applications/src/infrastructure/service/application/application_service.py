import logging
from asyncio import create_task, wait_for

from src.core.exc import AccessError, AlreadyExistsError, InternalError, NotFoundError
from src.domain.schemas import ApplicationMessage, ApplicationSchema, UserInfo, UserRole
from src.infrastructure.service.application.dependencies import (
    IApplicationProducer,
    IApplicationRepo,
    IMessageBox,
    IMessageSaver,
    IUnitOfWork,
)
from src.infrastructure.service.application.dto import application_dto


class ApplicationService:
    def __init__(
        self,
        uof: IUnitOfWork,
        application_repo: IApplicationRepo,
        repo_timeout: float,
        application_producer: IApplicationProducer,
        message_box: IMessageBox,
        message_saver: IMessageSaver,
    ) -> None:
        self.uof: IUnitOfWork = uof
        self.application_repo: IApplicationRepo = application_repo
        self.repo_timeout: float = repo_timeout

        self.application_producer: IApplicationProducer = application_producer
        self.message_box: IMessageBox[ApplicationSchema] = message_box
        self.message_saver: IMessageSaver = message_saver

        self.logger = logging.getLogger("ApplicationService")

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

        create_task(self.__put_application_in_box(application))

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

    async def __put_application_in_box(self, application: ApplicationSchema) -> None:
        try:
            self.message_box.add_message(application)
        except OverflowError:
            messages = self.message_box.get_messages()
            self.message_box.drop_box()
            self.message_box.add_message(application)

            await self.__publish_applications(messages)

    async def __publish_applications(self, messages: list[ApplicationSchema | None]) -> None:
        amounts = count_applications_amount(messages)
        applications = [ApplicationMessage(vacancy_id=vacancy_id, amount=amount) for vacancy_id, amount in amounts.items()]

        try:
            await self.application_producer.publish_applications(applications)
        except InternalError as e:
            self.logger.error(f"can't send application messages in broker, details={e}")
            create_task(self.__save_applications(applications))

    async def __save_applications(self, applications: list[ApplicationMessage]) -> None:
        try:
            await self.message_saver.save_messages(applications)
        except InternalError as e:
            self.logger.critical(f"can't save application messages in message saver, details={e}")


def count_applications_amount(messages: list[ApplicationSchema | None]) -> dict[int, int]:
    amounts = {}
    for message in messages:
        if message is not None:
            amounts[message.vacancy_id] = amounts.get(message.vacancy_id, 0) + 1

    return amounts
