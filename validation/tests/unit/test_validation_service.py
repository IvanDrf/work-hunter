from asyncio import gather
from itertools import chain

from pytest import mark, raises

from src.core.exc import InternalError
from src.infrastructure.service import ValidationService


@mark.asyncio
async def test_validation_metro(validation_service: ValidationService, cities_and_metro: dict[str, set[str]]) -> None:
    for city, stations in cities_and_metro.items():
        for station in stations:
            valid = (
                (city, station),
                (city.capitalize(), station),
                (city, station.capitalize()),
                (city.capitalize(), station.capitalize()),
                (city.lower(), station),
                (city, station.lower()),
                (city.lower(), station.lower()),
            )

            invalid = (
                (station, city),
                (station.capitalize(), city),
                (station.capitalize(), city.capitalize()),
                (station.lower(), city.lower()),
                (station.lower(), city),
                (station.lower(), city.lower()),
                (city, 2 * station),
            )

            results = await gather(*[validation_service.is_metro_valid(city, station) for city, station in chain(valid, invalid)])

            for res in results[: len(valid)]:
                assert res is True

            for res in results[len(valid) :]:
                assert res is False


@mark.asyncio
async def test_validation_metro_slow_service(
    validation_slow_service: ValidationService,
    cities_and_metro: dict[str, set[str]],
) -> None:
    async def assert_exc(city: str, station: str) -> None:
        with raises(InternalError) as e:
            await validation_slow_service.is_metro_valid(city, station)
            e.match(f"can't check is {station=} exists in {city=}")

    for city, stations in cities_and_metro.items():
        await gather(*[assert_exc(city, station) for station in stations])


@mark.asyncio
async def test_validation_city(validation_service: ValidationService, cities: set[str]) -> None:
    for city in cities:
        valid = (city, city.capitalize(), city.lower())
        invalid = (city[::-1], city[::-1].capitalize(), city[::-1].lower())

        results = await gather(*[validation_service.is_city_valid(city) for city in chain(valid, invalid)])

        for res in results[: len(valid)]:
            assert res is True

        for res in results[len(valid) :]:
            assert res is False


@mark.asyncio
async def test_validation_city_slow_service(
    validation_slow_service: ValidationService,
    cities: set[str],
) -> None:
    async def assert_exc(city: str) -> None:
        with raises(InternalError) as e:
            await validation_slow_service.is_city_valid(city)
            e.match(f"can't check is {city=} exists")

        await gather(*[assert_exc(city) for city in cities])
