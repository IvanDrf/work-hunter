from typing import Protocol

from src.domain.schemas import Application


class IApplicationService(Protocol):
    async def increase_vacancy_applications(self, application: Application) -> None: ...
