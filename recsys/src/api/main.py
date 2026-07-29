from contextlib import asynccontextmanager

from fastapi import FastAPI, HTTPException, Path, Query

from src.config import settings
from src.embeddings.e5_model import E5Model
from src.models.search import SearchResponse
from src.models.vacancy import Vacancy
from src.pipeline.recommendation_pipeline import RecommendationPipeline
from src.preprocess.e5_small_formatter import E5Formatter
from src.vector_store.qdrant_store import QdrantStore

pipeline: RecommendationPipeline = None

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


@asynccontextmanager
async def lifespan(app: FastAPI):
    global pipeline
    print("Запуск API")

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

    print("Сервис готов к приему запросов!")

    yield

    print("Остановка сервиса")


app = FastAPI(
    title="Work Hunter RecSys API",
    description="API для поиска и рекомендаций вакансий",
    version="1.0.0",
    lifespan=lifespan,
)


@app.get("/")
def health_check():
    return {"status": "ok", "message": "Work Hunter API is running"}


@app.get("/api/search", response_model=SearchResponse)
def search_vacancies(
    query: str = Query(
        ..., description="Текстовый запрос (например: 'Ищу удаленку на Python')"
    ),
    limit: int = Query(5, ge=1, le=50, description="Сколько вакансий вернуть"),
):
    """
    Поиск вакансий по смыслу текста.
    """
    if not query.strip():
        raise HTTPException(
            status_code=400, detail="Поисковый запрос не может быть пустым."
        )

    return pipeline.search_by_query(query=query, limit=limit)


@app.get("/api/similar/{vacancy_id}", response_model=SearchResponse)
def get_similar_vacancies(
    vacancy_id: int = Path(..., gt=0, description="Уникальный ID вакансии в базе"),
    limit: int = Query(5, ge=1, le=50, description="Сколько похожих вакансий вернуть"),
):
    """
    Поиск вакансий, похожих на заданную по её ID.
    """
    try:
        return pipeline.search_by_id(vacancy_id=vacancy_id, limit=limit)
    except ValueError as e:
        raise HTTPException(status_code=404, detail=str(e))


@app.post("/api/vacancies", status_code=201)
def add_vacancy(vacancy: Vacancy):
    """
    Add a new vacancy to the database.
    """
    pipeline.upsert_vacancy(vacancy)

    return {
        "status": "success",
        "message": f"Vacancy {vacancy.vacancy_id} successfully added.",
    }


@app.put("/api/vacancies/{vacancy_id}")
def update_vacancy(vacancy_id: int, vacancy: Vacancy):
    """
    Update an existing vacancy.
    """
    if vacancy_id != vacancy.vacancy_id:
        raise HTTPException(
            status_code=400,
            detail="The vacancy_id in the URL does not match the vacancy_id in the request body.",
        )

    pipeline.upsert_vacancy(vacancy)

    return {
        "status": "success",
        "message": f"Vacancy {vacancy_id} successfully updated.",
    }


@app.delete("/api/vacancies/{vacancy_id}")
def delete_vacancy(vacancy_id: int):
    """
    Delete a vacancy by its ID.
    """
    pipeline.delete_vacancy(vacancy_id)

    return {
        "status": "success",
        "message": f"Vacancy {vacancy_id} successfully deleted.",
    }
