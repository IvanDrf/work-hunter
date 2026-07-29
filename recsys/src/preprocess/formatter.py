from abc import ABC, abstractmethod

from src.models.vacancy import Vacancy


class VacancyFormatter(ABC):
    """
    Абстрактный класс для форматирования вакансий в текст.
    """

    @property
    @abstractmethod
    def formatter_name(self) -> str:
        """Название форматтера для логирования и конфигов."""

    @abstractmethod
    def format(self, vacancy: Vacancy) -> str:
        """
        Принимает данные вакансии и возвращает строку для векторизации.
        """
