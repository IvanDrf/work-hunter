from abc import ABC, abstractmethod

from src.models.embeddings import VectorRepresentation


class EmbeddingModel(ABC):
    """
    Абстрактный класс для моделей векторизации.
    """

    @property
    @abstractmethod
    def vector_name(self) -> str:
        """Имя вектора для маршрутизации в БД."""

    @property
    @abstractmethod
    def dimension(self) -> int:
        """Размерность эмбеддинга модели."""

    @abstractmethod
    def encode_document(self, texts: str | list[str]) -> list[VectorRepresentation]:
        """
        Принимает текст или список текстов.
        """

    @abstractmethod
    def encode_query(self, text: str) -> VectorRepresentation:
        """
        Поисковый запрос векторизуется одиночно для search_by_query.
        """
