from typing import Any

from src.embeddings.base import EmbeddingModel
from src.models.search import SearchResponse
from src.models.vacancy import Vacancy
from src.preprocess.formatter import VacancyFormatter
from src.vector_store.base import VectorStore


class RecommendationPipeline:
    def __init__(
        self,
        formatter: VacancyFormatter,
        embedding_model: EmbeddingModel,
        vector_store: VectorStore,
        name: str = "default_pipeline",
    ):
        self.name = name
        self.formatter = formatter
        self.embedding_model = embedding_model
        self.vector_store = vector_store

    def upsert_vacancy(self, vacancy: Vacancy | list[Vacancy]) -> bool:
        """
        Upserts a vacancy or a list of vacancies.
        """
        vacancies = [vacancy] if isinstance(vacancy, Vacancy) else vacancy

        if not vacancies:
            return True

        formatted_texts = [self.formatter.format(v) for v in vacancies]
        representations = self.embedding_model.encode_document(formatted_texts)

        ids = [v.vacancy_id for v in vacancies]
        metadatas = [
            v.model_dump(
                exclude={"description", "requirements", "conditions"}, exclude_none=True
            )
            for v in vacancies
        ]

        return self.vector_store.upsert(
            ids=ids,
            representations=representations,
            metadatas=metadatas,
        )

    def delete_vacancy(self, vacancy_id: int | list[int]) -> bool:
        """
        Удаляет вакансию или список вакансий по их ID.
        """
        ids = [vacancy_id] if isinstance(vacancy_id, int) else vacancy_id

        if not ids:
            return True

        return self.vector_store.delete(ids=ids)

    def search_by_query(
        self,
        query: str,
        limit: int = 5,
        offset: int = 0,
        filters: dict[str, Any] | None = None,
    ) -> SearchResponse:
        """
        Поиск вакансий по текстовому поисковому запросу.
        """
        if not query.strip():
            return SearchResponse(items=[])

        query_representation = self.embedding_model.encode_query(query)

        return self.vector_store.search_similar(
            query_representation=query_representation,
            limit=limit,
            offset=offset,
            filters=filters,
        )

    def search_by_id(
        self,
        vacancy_id: int,
        limit: int = 5,
        offset: int = 0,
        filters: dict[str, Any] | None = None,
    ) -> SearchResponse:
        """
        Поиск похожих на вакансию из базы.
        """
        record = self.vector_store.get_by_id(vacancy_id, with_vector=True)
        if not record or not record.representation:
            raise ValueError(f"Вакансия с ID {vacancy_id} не найдена.")

        return self.vector_store.search_similar(
            query_representation=record.representation,
            limit=limit,
            offset=offset,
            filters=filters,
            exclude_id=vacancy_id,
        )
