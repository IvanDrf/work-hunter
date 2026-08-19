from logging import CRITICAL, DEBUG, ERROR, INFO, WARNING
from typing import Final

from pytest import mark, raises
from sqlalchemy.exc import SQLAlchemyError

from src.core.exc import AccessError, ArgumentError, InternalError, NotFoundError
from src.utils.catch_error import catch_raise_error

levels = {
    DEBUG: "DEBUG",
    INFO: "INFO",
    WARNING: "WARNING",
    ERROR: "ERROR",
    CRITICAL: "CRITICAL",
}


@mark.parametrize(
    "catch_err, raise_err, message",
    [
        (SQLAlchemyError, InternalError, "catch_alchemy_internal"),
        (InternalError, AccessError, "catch_internal_access"),
        (NotFoundError, ArgumentError, "catch_not_found_argument"),
    ],
)
@mark.asyncio
async def test_catch_raise_error(catch_err, raise_err, message, caplog) -> None:
    SUCCESS: Final[str] = "success"

    for level, level_str in levels.items():
        caplog.set_level(level)

        @catch_raise_error((catch_err,), raise_error=raise_err, logger_level=level_str.lower(), message=message)  # type: ignore
        async def wrap(err: type[Exception], *, raise_err: bool, message: str) -> str:
            if raise_err is True:
                raise err(message)

            return SUCCESS

        with raises(raise_err):
            await wrap(catch_err, raise_err=True, message=message)

        for record in caplog.records:
            assert record.levelname == level_str
            assert record.getMessage() == f"wrap: {message}"

        caplog.clear()

        res = await wrap(catch_err, raise_err=False, message=message)
        assert res == SUCCESS
        assert len(caplog.records) == 0
