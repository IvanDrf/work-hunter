from sqlalchemy import BIGINT, VARCHAR
from sqlalchemy.orm import Mapped, mapped_column, relationship

from src.domain.models.base import Base


class TagORM(Base):
    __tablename__ = 'tags'

    tag_id: Mapped[int] = mapped_column(
        BIGINT, primary_key=True, autoincrement=True, index=True
    )

    tag: Mapped[str] = mapped_column(
        VARCHAR(40), nullable=False, unique=True, index=True
    )

    vacancies: Mapped[list['VacancyORM']] = relationship(  # type: ignore
        back_populates='tags', secondary='vacancies_to_tags',
        cascade='save-update, merge, delete'
    )
