import logging
from typing import Annotated

from fastapi import APIRouter, Depends, HTTPException, Request, status
from slowapi import Limiter
from slowapi.util import get_remote_address

from src.api.dependencies import IValidationService, get_validation_service
from src.api.middleware import api_key_middleware
from src.core.exc import InternalError
from src.domain.schemas import CitySchema, MetroSchema

validation_router = APIRouter(prefix="/api", dependencies=[Depends(api_key_middleware)])

limiter = Limiter(key_func=get_remote_address)


@validation_router.post("/metro", status_code=status.HTTP_204_NO_CONTENT)
@limiter.limit("30/second")
async def validate_metro(
    payload: MetroSchema,
    validation_service: Annotated[IValidationService, Depends(get_validation_service)],
    request: Request,
) -> None:
    try:
        if not await validation_service.is_metro_valid(payload.city, payload.metro):
            raise HTTPException(status_code=status.HTTP_404_NOT_FOUND)

    except InternalError as e:
        logging.critical(f"can't validate metro and city, details={e}")
        raise HTTPException(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail={"error": str(e)})


@validation_router.post("/city", status_code=status.HTTP_204_NO_CONTENT)
@limiter.limit("30/second")
async def validate_city(
    payload: CitySchema,
    validation_service: Annotated[IValidationService, Depends(get_validation_service)],
    request: Request,
):
    try:
        if not await validation_service.is_city_valid(payload.city):
            raise HTTPException(status_code=status.HTTP_404_NOT_FOUND)

    except InternalError as e:
        logging.critical(f"can't validate city, details={e}")
        raise HTTPException(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail={"error": str(e)})
