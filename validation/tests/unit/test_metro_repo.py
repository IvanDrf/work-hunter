from pytest import mark

from src.infrastructure.repo import MetroRedisRepo
from hypothesis import strategies as st
from hypothesis import given


@mark.asyncio
async def test_metro_redis_repo(metro_repo: MetroRedisRepo, city_and_metro: tuple[str, list[str]]) -> None:
    for city, stations in city_and_metro:
        city = city.lower()
        for station in stations:
            station = station.lower()

            assert await metro_repo.is_metro_exists(city, station) is True


@mark.asyncio
@given(city=st.text(), metro=st.text())
async def test_metro_redis_not_exists(metro_repo: MetroRedisRepo, city: str, metro: str) -> None:
    assert await metro_repo.is_metro_exists(city, metro) is False
