from uuid import UUID

from pydantic import BaseModel, Field

from src.domain.schemas.user import UserInfo


class ApplicationSchema(BaseModel):
    vacancy_id: int = Field(ge=0)
    user_info: UserInfo


class ApplicationMessage(BaseModel):
    message_id: UUID
    vacancy_id: int = Field(ge=0)
