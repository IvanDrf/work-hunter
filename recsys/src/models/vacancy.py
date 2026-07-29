from pydantic import BaseModel


class Vacancy(BaseModel):
    vacancy_id: int
    title: str
    author_name: str
    author_id: str | None = None

    # Текстовые поля для описания
    description: str | None = None
    requirements: str | None = None
    conditions: str | None = None
    tags: str | None = None

    # Фильтры
    city: str | None = None
    metro: str | None = None
    currency: str | None = None
    remote_type: str | None = None
    time_type: str | None = None

    # Числовые фильтры
    salary_min: int | None = None
    salary_max: int | None = None
    experience_min: int | None = None
    experience_max: int | None = None
