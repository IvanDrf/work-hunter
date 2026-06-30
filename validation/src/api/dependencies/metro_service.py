from typing import Protocol


class IMetroService(Protocol):
    async def is_metro_valid(self, city: str, metro: str) -> bool: ...
    async def close(self) -> None: ...


metro_service: IMetroService | None = None


def init_metro_service(service: IMetroService) -> None:
    global metro_service
    metro_service = service


def get_metro_service() -> IMetroService:
    if metro_service is None:
        raise RuntimeError("metro service is None")

    return metro_service
