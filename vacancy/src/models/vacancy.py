from datetime import datetime
from enum import Enum as PyEnum

from sqlalchemy import BIGINT, INT, TIMESTAMP, VARCHAR, CheckConstraint, Enum, Text
from sqlalchemy.orm import Mapped, mapped_column, relationship

from src.models.base import Base


class RemoteType(PyEnum):
    OFFICE = 0
    REMOTE = 1
    HYBRID = 2
    ANY = 3


class TimeType(PyEnum):
    FULL = 0
    PART = 1


class Currency(PyEnum):
    RUB = 0
    USD = 1
    EUR = 2


class VacancyStatus(PyEnum):
    MODERATING = 0
    PUBLISHED = 1
    UPDATED = 2
    CLOSED = 3
    DELETED = 4


class VacancyORM(Base):
    __tablename__ = 'vacancies'

    vacancy_id: Mapped[int] = mapped_column(
        BIGINT, primary_key=True, autoincrement=True, index=True
    )

    title: Mapped[str] = mapped_column(VARCHAR(150), nullable=False)
    description: Mapped[str] = mapped_column(Text, nullable=False)

    requirements: Mapped[str] = mapped_column(Text, nullable=False)
    conditions: Mapped[str] = mapped_column(Text, nullable=False)

    salary_min: Mapped[int] = mapped_column(
        INT, CheckConstraint('salary_min >= 0', name='check_positive_salary_min'), nullable=False,
    )
    salary_max: Mapped[int] = mapped_column(
        INT, CheckConstraint('salary_max >= 0', name='check_positive_salary_max'), nullable=False,
    )

    currency: Mapped[Currency] = mapped_column(Enum(Currency), nullable=False)

    city: Mapped[str] = mapped_column(VARCHAR(150), nullable=True)
    metro: Mapped[str] = mapped_column(VARCHAR(100), nullable=True)

    remote_type: Mapped[RemoteType] = mapped_column(
        Enum(RemoteType), nullable=False
    )
    time_type: Mapped[TimeType] = mapped_column(Enum(TimeType), nullable=False)

    experience_min: Mapped[int] = mapped_column(
        INT, CheckConstraint('experience_min >= 0', name='check_non_negative_min_exp'), nullable=True
    )
    experience_max: Mapped[int] = mapped_column(
        INT, CheckConstraint('experience_max >= 0', name='check_non_negative_max_exp'), nullable=True
    )

    created_at: Mapped[datetime] = mapped_column(TIMESTAMP, nullable=False)
    updated_at: Mapped[datetime] = mapped_column(TIMESTAMP, nullable=True)
    published_at: Mapped[datetime] = mapped_column(TIMESTAMP, nullable=True)
    closed_at: Mapped[datetime] = mapped_column(TIMESTAMP, nullable=True)

    status: Mapped[VacancyStatus] = mapped_column(
        Enum(VacancyStatus), nullable=False)
    moderated_at: Mapped[datetime] = mapped_column(TIMESTAMP, nullable=True)
    moderator_comments: Mapped[str] = mapped_column(Text, nullable=True)

    views: Mapped[int] = mapped_column(
        BIGINT, CheckConstraint('views >= 0', name='check_non_negative_views'), nullable=True
    )

    applications_count: Mapped[int] = mapped_column(
        BIGINT, CheckConstraint('applications_count >= 0', name='check_non_negative_applications'), nullable=False, default=0
    )

    tags = relationship(
        'TagORM', secondary='vacancies_to_tags', back_populates='vacancies'
    )
