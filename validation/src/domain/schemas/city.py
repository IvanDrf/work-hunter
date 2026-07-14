from pydantic import BaseModel, Field

MAX_CITY_NAME_LENGTH = 200


class CitySchema(BaseModel):
    city: str = Field(min_length=1, max_length=MAX_CITY_NAME_LENGTH)
