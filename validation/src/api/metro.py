from fastapi import APIRouter, Depends

from src.api.middleware import api_key_middleware

metro_router = APIRouter(prefix="/api/metro", dependencies=[Depends(api_key_middleware)])
