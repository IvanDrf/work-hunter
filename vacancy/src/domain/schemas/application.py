from uuid import UUID

from pydantic import BaseModel, Field


class ApplicationMessage(BaseModel):
    message_id: UUID
    vacancy_id: int = Field(ge=0)
