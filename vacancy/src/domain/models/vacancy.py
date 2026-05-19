from datetime import datetime
from typing import TypeAlias
from uuid import UUID as PyUUID

from sqlalchemy import BIGINT, INT, TIMESTAMP, UUID, VARCHAR, CheckConstraint, Enum, Text
from sqlalchemy.orm import Mapped, mapped_column, relationship

from src.domain.models.base import Base
from src.domain.types.enums import Currency, RemoteType, TimeType, VacancyStatus

Year: TypeAlias = int
Money: TypeAlias = int


class VacancyORM(Base):
    __tablename__ = "vacancies"

    vacancy_id: Mapped[int] = mapped_column(BIGINT, primary_key=True, autoincrement=True, index=True)
    author_id: Mapped[PyUUID] = mapped_column(UUID, index=True, nullable=False)
    author_name: Mapped[str] = mapped_column(VARCHAR(75), nullable=False)

    title: Mapped[str] = mapped_column(VARCHAR(150), nullable=False)
    description: Mapped[str] = mapped_column(Text, nullable=False)

    requirements: Mapped[str] = mapped_column(Text, nullable=False)
    conditions: Mapped[str] = mapped_column(Text, nullable=False)

    salary_min: Mapped[Money | None] = mapped_column(
        INT, CheckConstraint("salary_min >= 0", name="check_positive_salary_min"), nullable=True, default=0
    )
    salary_max: Mapped[Money | None] = mapped_column(
        INT, CheckConstraint("salary_max >= 0", name="check_positive_salary_max"), nullable=True, default=0
    )

    currency: Mapped[Currency] = mapped_column(Enum(Currency), nullable=False)

    city: Mapped[str | None] = mapped_column(VARCHAR(150), nullable=True)
    metro: Mapped[str | None] = mapped_column(VARCHAR(100), nullable=True)

    remote_type: Mapped[RemoteType] = mapped_column(Enum(RemoteType), nullable=False)
    time_type: Mapped[TimeType] = mapped_column(Enum(TimeType), nullable=False)

    experience_min: Mapped[Year | None] = mapped_column(
        INT,
        CheckConstraint("experience_min >= 0", name="check_non_negative_min_exp"),
        nullable=True,
        default=None,
    )
    experience_max: Mapped[Year | None] = mapped_column(
        INT,
        CheckConstraint("experience_max >= 0", name="check_non_negative_max_exp"),
        nullable=True,
        default=None,
    )

    created_at: Mapped[datetime] = mapped_column(TIMESTAMP(timezone=True), nullable=False)
    updated_at: Mapped[datetime | None] = mapped_column(TIMESTAMP(timezone=True), nullable=True)
    published_at: Mapped[datetime | None] = mapped_column(TIMESTAMP(timezone=True), nullable=True)
    closed_at: Mapped[datetime | None] = mapped_column(TIMESTAMP(timezone=True), nullable=True)
    deleted_at: Mapped[datetime | None] = mapped_column(TIMESTAMP(timezone=True), nullable=True)

    status: Mapped[VacancyStatus] = mapped_column(Enum(VacancyStatus), nullable=False)
    moderated_at: Mapped[datetime | None] = mapped_column(TIMESTAMP(timezone=True), nullable=True)
    moderator_comments: Mapped[str | None] = mapped_column(Text, nullable=True)

    views: Mapped[int] = mapped_column(
        BIGINT, CheckConstraint("views >= 0", name="check_non_negative_views"), nullable=False, default=0
    )

    applications_count: Mapped[int] = mapped_column(
        BIGINT, CheckConstraint("applications_count >= 0", name="check_non_negative_applications"), nullable=False, default=0
    )

    tags: Mapped[list["TagORM"]] = relationship(  # noqa # type: ignore
        back_populates="vacancies", secondary="vacancies_to_tags", cascade="save-update, merge, delete"
    )
