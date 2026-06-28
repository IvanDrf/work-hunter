from typing import Protocol

from src.domain.schemas import ApplicationMessage


class IApplicationService(Protocol):
    async def increase_vacancy_applications(self, application: ApplicationMessage) -> None: ...
