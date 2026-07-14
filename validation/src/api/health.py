from fastapi import APIRouter, status

health_router = APIRouter()


@health_router.get("/health", status_code=status.HTTP_200_OK)
async def health() -> dict[str, str]:
    return {"status": "AVAILABLE"}
