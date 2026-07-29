from src.models.vacancy import Vacancy
from src.preprocess.formatter import VacancyFormatter


class E5Formatter(VacancyFormatter):
    """
    Форматтер для модели multilingual-e5-small.
    """

    @property
    def formatter_name(self) -> str:
        return "multilingual-e5-small-formatter"

    def format(self, vacancy: Vacancy) -> str:
        parts = [f"Должность: {vacancy.title}", f"Компания: {vacancy.author_name}"]

        if vacancy.requirements:
            parts.append(f"Требования: {vacancy.requirements}")

        if vacancy.tags:
            parts.append(f"Ключевые слова: {vacancy.tags}")

        return ". ".join(parts) + "."
