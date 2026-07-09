from string import ascii_letters

from hypothesis import given
from hypothesis import strategies as st
from pydantic import ValidationError
from pytest import raises

from src.domain.schemas import CitySchema
from src.domain.schemas.city import MAX_CITY_NAME_LENGTH


@given(city=st.text(alphabet=ascii_letters, min_size=1, max_size=MAX_CITY_NAME_LENGTH))
def test_validate_metro_schema(city: str) -> None:
    schema = CitySchema(city=city)

    assert schema.city == city


@given(
    short_city=st.text(alphabet=ascii_letters, max_size=0),
    long_city=st.text(alphabet=ascii_letters, min_size=MAX_CITY_NAME_LENGTH + 1, max_size=2 * MAX_CITY_NAME_LENGTH),
)
def test_validate_metro_schema_invalid_city(short_city: str, long_city: str) -> None:
    with raises(ValidationError):
        CitySchema(city=short_city)

    with raises(ValidationError):
        CitySchema(city=long_city)
