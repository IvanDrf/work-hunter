from typing import Any

from pydantic import BaseModel

from src.models.embeddings import VectorRepresentation


class SearchResultItem(BaseModel):
    id: int
    score: float
    metadata: dict[str, Any]


class SearchResponse(BaseModel):
    items: list[SearchResultItem]


class VectorRecord(BaseModel):
    id: int
    representation: VectorRepresentation | None = None
    metadata: dict[str, Any]
