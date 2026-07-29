from sentence_transformers import SentenceTransformer

from src.embeddings.base import EmbeddingModel
from src.models.embeddings import VectorRepresentation


class E5Model(EmbeddingModel):
    def __init__(
        self, model_name: str = "intfloat/multilingual-e5-small", batch_size: int = 32
    ):
        self.model = SentenceTransformer(model_name)
        self.batch_size = batch_size
        self._vector_name = ""

    @property
    def vector_name(self) -> str:
        return self._vector_name

    @property
    def dimension(self) -> int:
        return self.model.get_embedding_dimension()

    def _encode_batch(self, texts: list[str]) -> list[VectorRepresentation]:
        embeddings = self.model.encode(
            texts,
            batch_size=self.batch_size,
            show_progress_bar=False,
            convert_to_numpy=True,
        )
        return [
            VectorRepresentation(dense=emb.tolist(), sparse=None) for emb in embeddings
        ]

    def encode_document(self, texts: str | list[str]) -> list[VectorRepresentation]:
        text_list = [texts] if isinstance(texts, str) else texts

        prefixed_texts = [f"passage: {t}" for t in text_list]

        return self._encode_batch(prefixed_texts)

    def encode_query(self, text: str) -> VectorRepresentation:
        prefixed_text = f"query: {text}"
        results = self._encode_batch([prefixed_text])
        return results[0]
