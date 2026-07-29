from abc import ABC, abstractmethod
from typing import Any

from src.models.embeddings import VectorRepresentation
from src.models.search import SearchResponse, VectorRecord


class VectorStore(ABC):
    @abstractmethod
    def upsert(
        self,
        ids: list[int],
        representations: list[VectorRepresentation],
        metadatas: list[dict[str, Any]],
    ) -> bool:
        """
        Upserts a vacancy or a list of vacancies.
        """

    @abstractmethod
    def get_by_id(self, id: int, with_vector: bool = True) -> VectorRecord | None:
        """
        Достает одну запись из базы по ID.
        """

    @abstractmethod
    def delete(self, ids: list[int]) -> bool:
        """
        Удаляет записи по списку ID.
        """

    @abstractmethod
    def search_similar(
        self,
        query_representation: VectorRepresentation,
        limit: int = 5,
        offset: int = 0,
        filters: Any | None = None,
        exclude_id: int | None = None,
    ) -> SearchResponse:
        """
        Ищет ближайшие объекты по векторному представлению.
        """
