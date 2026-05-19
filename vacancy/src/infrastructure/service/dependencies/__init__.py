from src.infrastructure.service.dependencies.cache import ICache
from src.infrastructure.service.dependencies.repo import ITagRepo, IVacancyRepo
from src.infrastructure.service.dependencies.unit_of_work import IUnitOfWork

__all__ = [
    "ICache",
    "IVacancyRepo",
    "ITagRepo",
    "IUnitOfWork",
]
