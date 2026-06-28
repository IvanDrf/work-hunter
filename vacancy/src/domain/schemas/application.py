from pydantic import BaseModel, Field

from src.domain.schemas.user import UserInfo


class ApplicationMessage(BaseModel):
    vacancy_id: int = Field(ge=0)
    applications: int = Field(gt=0)

    user_info: UserInfo
