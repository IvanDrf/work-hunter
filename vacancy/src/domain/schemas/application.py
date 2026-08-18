from uuid import UUID

from pydantic import BaseModel, Field


class ApplicationMessage(BaseModel):
    message_id: UUID
    vacancy_id: int = Field(ge=0)

    def __hash__(self) -> int:
        return hash(self.message_id)

    def __eq__(self, value: object) -> bool:
        return isinstance(value, ApplicationMessage) and self.message_id == value.message_id
