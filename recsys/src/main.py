from src.config import settings
from src.embeddings.e5_model import E5Model
from src.pipeline.recommendation_pipeline import RecommendationPipeline
from src.preprocess.e5_small_formatter import E5Formatter
from src.vector_store.qdrant_store import QdrantStore

model = E5Model(model_name=settings.embedding_model, batch_size=64)

store = QdrantStore(
    host=settings.qdrant_host,
    port=settings.qdrant_port,
    collection_name=settings.collection_name,
    vector_size=model.dimension,
)

formatter = E5Formatter()

pipeline = RecommendationPipeline(
    formatter=formatter,
    embedding_model=model,
    vector_store=store,
)
