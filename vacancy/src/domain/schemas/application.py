from pydantic import BaseModel, Field


class Application(BaseModel):
    vacancy_id: int = Field(ge=0)
    applications: int = Field(gt=0)
