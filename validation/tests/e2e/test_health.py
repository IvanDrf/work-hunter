from fastapi import status
from fastapi.testclient import TestClient

from src.app.main import app


def test_health(requests_amount: int) -> None:
    for _ in range(requests_amount):
        with TestClient(app=app) as client:
            response = client.get(url="/health")
            assert response.status_code == status.HTTP_200_OK
            assert response.json() == {"status": "AVAILABLE"}

            response.close()
