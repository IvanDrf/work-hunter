from concurrent.futures import ThreadPoolExecutor

from fastapi import status
from fastapi.testclient import TestClient

from src.api.validation import limiter
from src.app.main import app
from src.core.config import CONFIG


def test_validate_city(cities_json) -> None:
    with TestClient(app=app) as client:
        cities_batch = []
        size = 15
        for city in cities_json:
            cities_batch.append(city)

            if len(cities_batch) >= size:
                send_requests_for_batch(client, cities_batch)
                cities_batch = []

        if cities_batch:
            send_requests_for_batch(client, cities_batch)


def test_validate_city_limiter(cities_json) -> None:
    RPS = 30
    requests_amount = 0
    limiter._storage.reset()

    with TestClient(app=app) as client:
        for city in cities_json:
            response = client.post(
                url="/api/city",
                headers={"api-key": CONFIG.api_key},
                json={"city": city},
            )

            requests_amount += 1
            if requests_amount > RPS:
                assert response.status_code == status.HTTP_429_TOO_MANY_REQUESTS
            else:
                assert response.status_code == status.HTTP_204_NO_CONTENT

            if requests_amount > 3 * RPS:  # limiter timout
                return


def send_requests_for_batch(client: TestClient, cities: list[str]) -> None:
    funcs = (
        assert_valid_request,
        assert_request_without_api_key,
        assert_request_invalid_body,
        assert_request_not_found_city,
    )

    with ThreadPoolExecutor(max_workers=15) as executor:
        limiter._storage.reset()  # no 429 error
        for city in cities:
            args = ((client, city), (client, city), (client, city), (client, 2 * city))

            [executor.submit(func, *arg) for func, arg in zip(funcs, args)]


def assert_valid_request(client: TestClient, city: str) -> None:
    response = client.post(
        url="/api/city",
        headers={"api-key": CONFIG.api_key},
        json={"city": city},
    )

    assert response.status_code == status.HTTP_204_NO_CONTENT
    response.close()


def assert_request_without_api_key(client: TestClient, city: str) -> None:
    response = client.post(
        url="/api/city",
        json={"city": city},
    )

    assert response.status_code == status.HTTP_401_UNAUTHORIZED
    response.close()


def assert_request_invalid_body(client: TestClient, city: str) -> None:
    response = client.post(
        url="/api/city",
        headers={"api-key": CONFIG.api_key},
        json={"city_": city},
    )

    assert response.status_code == status.HTTP_422_UNPROCESSABLE_CONTENT
    response.close()


def assert_request_not_found_city(client: TestClient, city: str) -> None:
    response = client.post(
        url="/api/city",
        headers={"api-key": CONFIG.api_key},
        json={"city": city},
    )

    assert response.status_code == status.HTTP_404_NOT_FOUND
    response.close()
