from typing import Protocol

from src.domain.schemas import ApplicationSchema, UserInfo


class IApplicationService(Protocol):
    async def update_application(self, application: ApplicationSchema) -> None: ...
    async def find_vacancies_ids_by_applications(self, user_info: UserInfo, *, limit: int, offset: int) -> list[int]: ...

    async def stop(self) -> None: ...
