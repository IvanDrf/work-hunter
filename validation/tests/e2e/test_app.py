from fastapi import status
from fastapi.testclient import TestClient

from src.api.validation import limiter
from src.app.main import app
from src.core.config import CONFIG


def test_health(requests_amount: int) -> None:
    for _ in range(requests_amount):
        with TestClient(app=app) as client:
            response = client.get(url="/health")
            assert response.status_code == status.HTTP_200_OK
            assert response.json() == {"status": "AVAILABLE"}

            response.close()


def test_validate_metro(cities_and_metro_json) -> None:
    with TestClient(app=app) as client:
        for city, stations in cities_and_metro_json:
            for station in stations:
                limiter._storage.reset()  # no 429 error

                assert_valid_request(client, city, station)
                assert_request_without_api_key(client, city, station)
                assert_request_invalid_body(client, city, station)
                assert_request_not_found_metro(client, city, 2 * station)


def test_validate_metro_limiter(cities_and_metro_json) -> None:
    RPS = 30
    requests_amount = 0
    limiter._storage.reset()

    with TestClient(app=app) as client:
        for city, stations in cities_and_metro_json:
            for station in stations:
                response = client.post(
                    url="/api/metro",
                    headers={"api-key": CONFIG.api_key, "Content-type": "application/json"},
                    json={"city": city, "metro": station},
                )

                requests_amount += 1
                if requests_amount > RPS:
                    assert response.status_code == status.HTTP_429_TOO_MANY_REQUESTS
                else:
                    assert response.status_code == status.HTTP_204_NO_CONTENT

                if requests_amount > 3 * RPS:
                    return


def assert_valid_request(client: TestClient, city: str, station: str) -> None:
    response = client.post(
        url="/api/metro",
        headers={"api-key": CONFIG.api_key, "Content-type": "application/json"},
        json={"city": city, "metro": station},
    )

    assert response.status_code == status.HTTP_204_NO_CONTENT
    response.close()


def assert_request_without_api_key(client: TestClient, city: str, station: str) -> None:
    response = client.post(
        url="/api/metro",
        headers={"Content-type": "application/json"},
        json={"city": city, "metro": station},
    )

    assert response.status_code == status.HTTP_401_UNAUTHORIZED
    response.close()


def assert_request_invalid_body(client: TestClient, city: str, station: str) -> None:
    response = client.post(
        url="/api/metro",
        headers={"api-key": CONFIG.api_key, "Content-type": "application/json"},
        json={"city_": city, "metro_": station},
    )

    assert response.status_code == status.HTTP_422_UNPROCESSABLE_CONTENT
    response.close()


def assert_request_not_found_metro(client: TestClient, city: str, station: str) -> None:
    response = client.post(
        url="/api/metro",
        headers={"api-key": CONFIG.api_key, "Content-type": "application/json"},
        json={"city": city, "metro": station},
    )

    assert response.status_code == status.HTTP_404_NOT_FOUND
    response.close()
