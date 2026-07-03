from string import ascii_letters

from hypothesis import given
from hypothesis import strategies as st
from pydantic import ValidationError
from pytest import raises

from src.domain.schemas.metro import MAX_CITY_NAME_LENGTH, MAX_METRO_NAME_LENGTH, ValidateMetroSchema


@given(
    city=st.text(alphabet=ascii_letters, min_size=1, max_size=MAX_CITY_NAME_LENGTH),
    metro=st.text(alphabet=ascii_letters, min_size=1, max_size=MAX_METRO_NAME_LENGTH),
)
def test_validate_metro_schema(city: str, metro: str) -> None:
    schema = ValidateMetroSchema(city=city, metro=metro)

    assert schema.city == city
    assert schema.metro == metro


@given(
    short_city=st.text(alphabet=ascii_letters, max_size=0),
    long_city=st.text(alphabet=ascii_letters, min_size=MAX_CITY_NAME_LENGTH + 1, max_size=2 * MAX_CITY_NAME_LENGTH),
)
def test_validate_metro_schema_invalid_city(short_city: str, long_city: str) -> None:
    VALID_METRO = "metro"

    with raises(ValidationError):
        ValidateMetroSchema(city=short_city, metro=VALID_METRO)

    with raises(ValidationError):
        ValidateMetroSchema(city=long_city, metro=VALID_METRO)


@given(
    short_metro=st.text(alphabet=ascii_letters, max_size=0),
    long_metro=st.text(alphabet=ascii_letters, min_size=MAX_METRO_NAME_LENGTH + 1, max_size=2 * MAX_METRO_NAME_LENGTH),
)
def test_validate_metro_schema_invalid_metro(short_metro: str, long_metro: str) -> None:
    VALID_CITY = "city"

    with raises(ValidationError):
        ValidateMetroSchema(city=VALID_CITY, metro=short_metro)

    with raises(ValidationError):
        ValidateMetroSchema(city=VALID_CITY, metro=long_metro)
