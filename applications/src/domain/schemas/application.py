from pydantic import BaseModel, Field

from src.domain.schemas.user import UserInfo


class ApplicationSchema(BaseModel):
    vacancy_id: int = Field(ge=0)
    user_info: UserInfo
