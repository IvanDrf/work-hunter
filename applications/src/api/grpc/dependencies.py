from typing import Protocol

from src.domain.schemas import ApplicationSchema


class IApplicationService(Protocol):
    async def update_application(self, application: ApplicationSchema) -> None: ...
    async def stop(self) -> None: ...
