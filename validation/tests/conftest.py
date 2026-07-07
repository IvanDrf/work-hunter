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
                for city_name, metro_stations in city.items():
                    yield (city_name, metro_stations)

        yield data_yield()
