from asyncio import sleep

from fakeredis.aioredis import FakeRedis
from pytest import fixture
from pytest_asyncio import fixture as async_fixture

from src.core.config import Config
from src.infrastructure.repo import CityRedisRepo, MetroRedisRepo, load_cities, load_cities_and_metro
from src.infrastructure.service import ValidationService


@async_fixture(scope="session")
async def metro_repo() -> MetroRedisRepo:
    redis = FakeRedis()

    await load_cities_and_metro(Config(), redis)

    return MetroRedisRepo(redis)


@async_fixture(scope="session")
async def city_repo() -> CityRedisRepo:
    redis = FakeRedis()

    await load_cities(Config(), redis)

    return CityRedisRepo(redis)


@fixture(scope="package")
def cities_and_metro() -> dict[str, set[str]]:
    return {
        "city1": {"a", "b", "c", "d"},
        "city2": {"aa", "bb", "cc", "dd"},
        "city3": {"aaa", "bbb", "ccc", "ddd"},
    }


@fixture(scope="package")
def cities() -> set[str]:
    return {"Москва", "Владивосток", "Нижний новгород", "Новосибирск", "Воронеж"}


@fixture(scope="package")
def repo_timeout() -> float:
    return 0.1


@fixture(scope="package")
def metro_dict_repo(cities_and_metro: dict[str, set[str]]):
    class MetroDictRepo:
        def __init__(self) -> None:
            self.storage: dict[str, set[str]] = {city: stations for city, stations in cities_and_metro.items()}

        async def is_metro_exists(self, city: str, metro: str) -> bool:
            return city in self.storage and metro in self.storage[city]

        async def close(self) -> None:
            self.storage = {}

    return MetroDictRepo()


@fixture(scope="package")
def city_dict_repo(cities: set[str]):
    class CityDictRepo:
        def __init__(self) -> None:
            self.storage: set[str] = cities

        async def is_city_exists(self, city: str) -> bool:
            return city in self.storage

        async def close(self) -> None:
            self.storage = set()

    return CityDictRepo()


@fixture(scope="package")
def metro_slow_dict_repo(cities_and_metro: dict[str, set[str]], repo_timeout: float):
    class MetroDictRepo:
        def __init__(self) -> None:
            self.storage: dict[str, set[str]] = {city: stations for city, stations in cities_and_metro.items()}

        async def is_metro_exists(self, city: str, metro: str) -> bool:
            await sleep(2 * repo_timeout)
            return city in self.storage and metro in self.storage[city]

        async def close(self) -> None:
            self.storage = {}

    return MetroDictRepo()


@fixture(scope="package")
def city_slow_dict_repo(cities: set[str], repo_timeout: float):
    class CityDictRepo:
        def __init__(self) -> None:
            self.storage: set[str] = cities

        async def is_city_exists(self, city: str) -> bool:
            await sleep(2 * repo_timeout)
            return city in self.storage

        async def close(self) -> None:
            self.storage = set()

    return CityDictRepo()


@fixture(scope="package")
def metro_service(metro_dict_repo, city_dict_repo, repo_timeout: float) -> ValidationService:
    return ValidationService(metro_repo=metro_dict_repo, city_repo=city_dict_repo, repo_timeout=repo_timeout)


@fixture(scope="package")
def metro_slow_service(metro_slow_dict_repo, city_slow_dict_repo, repo_timeout: float) -> ValidationService:
    return ValidationService(metro_repo=metro_slow_dict_repo, city_repo=city_slow_dict_repo, repo_timeout=repo_timeout)
