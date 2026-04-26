from sqlalchemy import BIGINT, ForeignKey
from sqlalchemy.orm import Mapped, mapped_column

from src.domain.models.base import Base


class VacanciesTagsORM(Base):
    __tablename__ = 'vacancies_to_tags'

    vacancy_id: Mapped[int] = mapped_column(
        BIGINT, ForeignKey('vacancies.vacancy_id', ondelete='CASCADE'), primary_key=True
    )

    tag_id: Mapped[int] = mapped_column(
        BIGINT, ForeignKey('tags.tag_id', ondelete='CASCADE'), primary_key=True
    )
