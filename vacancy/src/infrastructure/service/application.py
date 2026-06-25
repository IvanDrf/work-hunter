from src.core.exc import AccessError, ArgumentError
from src.domain.rules.user import is_user_employee
from src.domain.schemas import ApplicationMessage
from src.infrastructure.service.dependencies import IUnitOfWork, IVacancyRepo


class ApplicationService:
    def __init__(self, vacancy_repo: IVacancyRepo, uof: IUnitOfWork) -> None:
        self.vacancy_repo: IVacancyRepo = vacancy_repo
        self.uof_factory: IUnitOfWork = uof

    async def increase_vacancy_applications(self, application: ApplicationMessage) -> None:
        if not is_user_employee(application.user_info):
            raise AccessError("only employee can apply for the vacacny")

        async with self.uof_factory as uof:
            vacancy = await self.vacancy_repo.find_vacancy_by_id(uof, application.vacancy_id)
            if vacancy is None:
                raise ArgumentError(f"can't find vacancy with vacancy_id={application.vacancy_id}")

            await self.vacancy_repo.update_vacancy(
                uof=uof,
                vacancy_id=application.vacancy_id,
                fields={"applications_count": vacancy.applications_count + 1},
            )
