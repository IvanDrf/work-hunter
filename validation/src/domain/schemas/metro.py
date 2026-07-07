from pydantic import Field

from src.domain.schemas.city import CitySchema

MAX_METRO_NAME_LENGTH = 80


class MetroSchema(CitySchema):
    metro: str = Field(min_length=1, max_length=MAX_METRO_NAME_LENGTH)
