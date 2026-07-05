import logging
from typing import Annotated

from fastapi import APIRouter, Depends, HTTPException, Request, status
from slowapi import Limiter
from slowapi.util import get_remote_address

from src.api.dependencies import IMetroService, get_metro_service
from src.api.middleware import api_key_middleware
from src.core.exc import InternalError
from src.domain.schemas import ValidateMetroSchema

metro_router = APIRouter(prefix="/api/metro", dependencies=[Depends(api_key_middleware)])

limiter = Limiter(key_func=get_remote_address)


@metro_router.post("", status_code=status.HTTP_204_NO_CONTENT)
@limiter.limit("30/second")
async def validate_metro(
    payload: ValidateMetroSchema,
    metro_service: Annotated[IMetroService, Depends(get_metro_service)],
    request: Request,
) -> None:
    try:
        if not await metro_service.is_metro_valid(payload.city, payload.metro):
            raise HTTPException(status_code=status.HTTP_404_NOT_FOUND)

    except InternalError as e:
        logging.critical(f"can't validate metro and city, details={e}")
        raise HTTPException(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail={"error": str(e)})
