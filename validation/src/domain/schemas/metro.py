from pydantic import BaseModel, Field, field_validator

from src.core.exc import ArgumentError

MAX_CITY_NAME_LENGTH = 200
MAX_METRO_NAME_LENGTH = 27


class ValidateMetroSchema(BaseModel):
    city: str = Field(min_length=1, max_length=MAX_CITY_NAME_LENGTH)
    metro: str = Field(min_length=1, max_length=MAX_METRO_NAME_LENGTH)

    @field_validator("city")
    def validate_city(cls, value: str) -> str:
        if not value.isalpha():
            raise ArgumentError(f"city must contain only alphabet, but city={value}")

        return value

    @field_validator("metro")
    def validate_metro(cls, value: str) -> str:
        if not value.isalpha():
            raise ArgumentError(f"metro must contain only alphabet, but metro={value}")

        return value
