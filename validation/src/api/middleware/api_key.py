from typing import Annotated, Final

from fastapi import Depends, HTTPException, status
from fastapi.security import APIKeyHeader

from src.core.config import CONFIG

API_KEY: Final[str] = CONFIG.api_key


def api_key_middleware(
    api_key: Annotated[str, Depends(APIKeyHeader(name="api-key", description="api-key is required for requests"))],
) -> None:
    if api_key != API_KEY:
        raise HTTPException(status_code=status.HTTP_403_FORBIDDEN, detail={"error": "api-key is invalid"})
