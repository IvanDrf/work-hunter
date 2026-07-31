from ijson import items
from pytest import fixture


@fixture(scope="session")
def requests_amount() -> int:
    return 5


@fixture(scope="function")
def cities_and_metro_json():
    METRO_FILE = "metro.json"

    with open(METRO_FILE, "r", encoding="utf-8") as metro_file:
        objects = items(metro_file, "cities.item")

        def data_yield():
            for city in objects:
                yield from city.items()

        yield data_yield()


@fixture(scope="function")
def cities_json():
    CITIES_FILE = "cities.json"

    with open(CITIES_FILE, "r", encoding="utf-8") as cities_file:
        cities = items(cities_file, "cities.item")

        def data_yield():
            yield from cities

        yield data_yield()
