from sqlalchemy import BIGINT, VARCHAR
from sqlalchemy.orm import Mapped, mapped_column, relationship

from src.models.base import Base


class TagORM(Base):
    __tablename__ = 'tags'

    tag_id: Mapped[int] = mapped_column(
        BIGINT, primary_key=True, autoincrement=True, index=True
    )

    tag: Mapped[str] = mapped_column(
        VARCHAR(40), nullable=False, unique=True, index=True
    )

    vacancies = relationship(
        'VacancyORM', secondary='vacancies_to_tags', back_populates='tags'
    )
