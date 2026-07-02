from asyncio import gather

from pytest import mark, raises

from src.core.exc import InternalError
from src.infrastructure.service.metro_service import MetroService


@mark.asyncio
async def test_metro_service(metro_service: MetroService, cities_and_metro: dict[str, set[str]]) -> None:
    for city, stations in cities_and_metro.items():
        for station in stations:
            assert await metro_service.is_metro_valid(city, station) is True
            assert await metro_service.is_metro_valid(station, city) is False
            assert await metro_service.is_metro_valid(city, 2 * station) is False


@mark.asyncio
async def test_metro_slow_service(metro_slow_service: MetroService, cities_and_metro: dict[str, set[str]]) -> None:
    async def assert_exc(city: str, station: str) -> None:
        with raises(InternalError) as e:
            await metro_slow_service.is_metro_valid(city, station)
            e.match(f"can't check is {station=} exists in {city=}")

    for city, stations in cities_and_metro.items():
        await gather(*[assert_exc(city, station) for station in stations])
