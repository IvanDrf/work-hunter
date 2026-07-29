from typing import Any

from qdrant_client import QdrantClient
from qdrant_client.models import (
    Distance,
    Filter,
    HasIdCondition,
    PointIdsList,
    PointStruct,
    SparseVector,
    VectorParams,
)

from src.models.embeddings import VectorRepresentation
from src.models.search import SearchResponse, SearchResultItem, VectorRecord
from src.vector_store.base import VectorStore


class QdrantStore(VectorStore):
    def __init__(
        self,
        collection_name: str,
        vector_size: int,
        host: str = "localhost",
        port: int = 6333,
        distance: Distance = Distance.COSINE,
    ):
        self.client = QdrantClient(host=host, port=port)
        self.collection_name = collection_name
        self.vector_size = vector_size
        self.distance = distance

        self._ensure_collection_exists()

    def _ensure_collection_exists(self):
        """Чистое создание коллекции без лишних разреженных векторов."""
        if not self.client.collection_exists(self.collection_name):
            self.client.create_collection(
                collection_name=self.collection_name,
                vectors_config=VectorParams(
                    size=self.vector_size,
                    distance=self.distance,
                ),
            )

    def _build_qdrant_vector(
        self, representation: VectorRepresentation
    ) -> dict[str, Any]:
        """Вспомогательный метод для упаковки векторов для Qdrant."""
        vector_dict = {}
        if representation.dense is not None:
            vector_dict[""] = representation.dense

        if representation.sparse is not None:
            vector_dict["keywords"] = SparseVector(
                indices=representation.sparse.indices,
                values=representation.sparse.values,
            )
        return vector_dict

    def upsert(
        self,
        ids: list[int],
        representations: list[VectorRepresentation],
        metadatas: list[dict[str, Any]],
    ) -> bool:
        if not ids:
            return True

        points = [
            PointStruct(
                id=v_id,
                vector=self._build_qdrant_vector(rep),
                payload=meta,
            )
            for v_id, rep, meta in zip(ids, representations, metadatas)
        ]

        self.client.upsert(collection_name=self.collection_name, points=points)
        return True

    def get_by_id(self, id: int, with_vector: bool = True) -> VectorRecord | None:
        results = self.client.retrieve(
            collection_name=self.collection_name,
            ids=[id],
            with_vectors=with_vector,
            with_payload=True,
        )
        if not results:
            return None

        hit = results[0]
        representation = None

        if with_vector and hit.vector:
            dense_vec = (
                hit.vector.get("") if isinstance(hit.vector, dict) else hit.vector
            )

            representation = VectorRepresentation(
                dense=dense_vec,
                sparse=None,
            )

        return VectorRecord(
            id=hit.id, representation=representation, metadata=hit.payload
        )

    def delete(self, ids: list[int]) -> bool:
        if not ids:
            return True

        response = self.client.delete(
            collection_name=self.collection_name,
            points_selector=PointIdsList(points=ids),
        )
        return response.status == "completed"

    def search_similar(
        self,
        query_representation: VectorRepresentation,
        limit: int = 5,
        offset: int = 0,
        filters: Any | None = None,
        exclude_id: int | None = None,
    ) -> SearchResponse:
        if not query_representation.dense:
            raise ValueError("Query representation must contain dense vector")

        final_filter = filters

        if exclude_id is not None:
            exclude_condition = HasIdCondition(has_id=[exclude_id])

            if final_filter is None:
                final_filter = Filter(must_not=[exclude_condition])
            else:
                final_filter = Filter(must=[filters], must_not=[exclude_condition])

        response = self.client.query_points(
            collection_name=self.collection_name,
            query=query_representation.dense,
            query_filter=final_filter,
            limit=limit,
            offset=offset,
            with_payload=True,
        )

        items = [
            SearchResultItem(id=hit.id, score=hit.score, metadata=hit.payload)
            for hit in response.points
        ]

        return SearchResponse(items=items)
