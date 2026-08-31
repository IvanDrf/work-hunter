from src.core.exc import InternalError
from src.domain.schemas import ApplicationMessage
from src.infrastructure.service.dependencies import IUnitOfWork, IVacancyRepo

MESSAGE_BATCH_SIZE = 20


class ApplicationService:
    def __init__(self, vacancy_repo: IVacancyRepo, uof: IUnitOfWork) -> None:
        self.vacancy_repo: IVacancyRepo = vacancy_repo
        self.uof: IUnitOfWork = uof
        self.batch = MessageBatch(MESSAGE_BATCH_SIZE)

    async def stop(self) -> None:
        await self.uof.stop()

    async def increase_vacancy_applications(self, message: ApplicationMessage) -> None:
        try:
            self.batch.add_application(message)  # trying to add new application message
            return
        except OverflowError:
            messages = self.batch.get_applications()  # if overflow get applications from batch, batch is empty now

        fields = create_update_applications_fields(messages)

        try:
            async with self.uof as uof:
                await self.vacancy_repo.update_vacancies(uof, fields)  # trying to update applications in database
                self.batch.add_application(message)  # add new incoming message in batch, no overflow, batch is empty
        except InternalError:
            # if we can't update applications in database (InternalError) we should save old applications from batch in batch again
            [self.batch.add_application(app) for app in messages]
            raise  # raise again cuz we didn't handle incoming message, batch is full again


class MessageBatch:
    def __init__(self, size: int) -> None:
        if size <= 0:
            raise ValueError("invalid size value in message batch, must be a positive")

        self.size = size
        self.batch: set[ApplicationMessage] = set()

    def add_application(self, application: ApplicationMessage) -> None:
        if len(self.batch) >= self.size:
            raise OverflowError("batch is full")

        self.batch.add(application)

    def get_applications(self) -> set[ApplicationMessage]:
        saved_apps = self.batch
        self.batch = set()

        return saved_apps


def create_update_applications_fields(messages: set[ApplicationMessage]) -> list[dict]:
    seen = set()
    vacancy_ids = {}

    for message in messages:
        if message.message_id in seen:
            continue

        seen.add(message.message_id)
        vacancy_ids[message.vacancy_id] = vacancy_ids.get(message.vacancy_id, 0) + 1

    return [{"vacancy_id": vacancy_id, "applications_count": amount} for vacancy_id, amount in vacancy_ids.items()]
