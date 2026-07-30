from uuid import UUID

from sqlalchemy import BIGINT, CheckConstraint
from sqlalchemy import UUID as SQLUUID
from sqlalchemy.orm import Mapped, mapped_column

from src.domain.models.base import Base


class ApplicationORM(Base):
    __tablename__ = "applications"

    vacancy_id: Mapped[int] = mapped_column(
        BIGINT,
        CheckConstraint("vacancy_id >= 0", name="check_non_negative_vacancy_id"),
        nullable=False,
        primary_key=True,
    )
    user_id: Mapped[UUID] = mapped_column(SQLUUID, nullable=False, primary_key=True)
