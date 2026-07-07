from typing import Protocol
from uuid import UUID

from src.domain.models import ApplicationORM
from src.infrastructure.service.application.dependencies.uof import IUnitOfWork


class IApplicationRepo(Protocol):
    async def add_application(self, uof: IUnitOfWork, application: ApplicationORM) -> None: ...

    async def find_application(self, uof: IUnitOfWork, vacancy_id: int, user_id: UUID) -> ApplicationORM | None: ...
