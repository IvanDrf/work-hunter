from pydantic import BaseModel, Field

MAX_CITY_NAME_LENGTH = 200
MAX_METRO_NAME_LENGTH = 80


class ValidateMetroSchema(BaseModel):
    city: str = Field(min_length=1, max_length=MAX_CITY_NAME_LENGTH)
    metro: str = Field(min_length=1, max_length=MAX_METRO_NAME_LENGTH)
