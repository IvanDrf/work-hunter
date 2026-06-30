from typing import Annotated

from fastapi import APIRouter, Depends, HTTPException, status

from src.api.dependencies import IMetroService, get_metro_service
from src.api.middleware import api_key_middleware
from src.domain.schemas import ValidateMetroSchema

metro_router = APIRouter(prefix="/api/metro", dependencies=[Depends(api_key_middleware)])


@metro_router.post("/", status_code=status.HTTP_200_OK)
async def validate_metro(payload: ValidateMetroSchema, metro_service: Annotated[IMetroService, get_metro_service]) -> None:
    if not metro_service.is_metro_valid(payload.city, payload.metro):
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND)
