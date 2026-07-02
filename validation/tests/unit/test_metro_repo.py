from hypothesis import given
from hypothesis import strategies as st
from pytest import mark

from src.infrastructure.repo import MetroRedisRepo


@mark.asyncio
async def test_metro_redis_repo(metro_repo: MetroRedisRepo, cities_and_metro_json: tuple[str, list[str]]) -> None:
    for city, stations in cities_and_metro_json:
        city = city.lower()
        for station in stations:
            station = station.lower()

            assert await metro_repo.is_metro_exists(city, station) is True


@mark.asyncio
@given(city=st.text(), metro=st.text())
async def test_metro_redis_not_exists(metro_repo: MetroRedisRepo, city: str, metro: str) -> None:
    assert await metro_repo.is_metro_exists(city, metro) is False
