from src.core.exc import ArgumentError
from src.domain.schemas import ApplicationMessage
from src.infrastructure.service.dependencies import IUnitOfWork, IVacancyRepo


class ApplicationService:
    def __init__(self, vacancy_repo: IVacancyRepo, uof: IUnitOfWork) -> None:
        self.vacancy_repo: IVacancyRepo = vacancy_repo
        self.uof_factory: IUnitOfWork = uof

    async def increase_vacancy_applications(self, message: ApplicationMessage) -> None:
        async with self.uof_factory as uof:
            vacancy = await self.vacancy_repo.find_vacancy_by_id(uof, message.vacancy_id)
            if vacancy is None:
                raise ArgumentError(f"can't find vacancy with vacancy_id={message.vacancy_id}")

            await self.vacancy_repo.update_vacancy(
                uof=uof,
                vacancy_id=message.vacancy_id,
                fields={"applications_count": vacancy.applications_count + 1},
            )
