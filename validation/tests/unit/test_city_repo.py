from collections.abc import Iterable

from hypothesis import given
from hypothesis import strategies as st
from pytest import mark

from src.infrastructure.repo import CityRedisRepo


@mark.asyncio
async def test_city_redis_repo(city_repo: CityRedisRepo, cities_json: Iterable[str]) -> None:
    for city in cities_json:
        assert await city_repo.is_city_exists(city.lower()) is True


@given(city=st.text())
@mark.asyncio
async def test_city_redis_repo_not_exists(city_repo: CityRedisRepo, city: str) -> None:
    assert await city_repo.is_city_exists(city) is False
