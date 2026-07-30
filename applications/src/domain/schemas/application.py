from pydantic import BaseModel, Field, TypeAdapter

from src.domain.schemas.user import UserInfo


class ApplicationSchema(BaseModel):
    vacancy_id: int = Field(ge=0)
    user_info: UserInfo


class ApplicationMessage(BaseModel):
    vacancy_id: int = Field(ge=0)
    amount: int = Field(ge=0)


Messages = TypeAdapter(list[ApplicationMessage])
