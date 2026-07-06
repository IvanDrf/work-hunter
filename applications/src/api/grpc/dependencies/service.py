from typing import Protocol

from src.domain.schemas import ApplicationSchema


class IApplicationService(Protocol):
    async def update_applications(self, application: ApplicationSchema) -> None: ...
