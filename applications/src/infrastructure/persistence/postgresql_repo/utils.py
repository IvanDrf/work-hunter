import logging
from typing import Literal

LoggerLevel = Literal["debug", "info", "warning", "error", "critical"]


def catch_and_raise(
    catch_exc: type[Exception] | tuple[type[Exception], ...],
    raise_exc: type[Exception],
    message: str,
    logger_level: LoggerLevel = "critical",
):
    def decorator(func):
        levels = {
            "debug": logging.debug,
            "info": logging.info,
            "warning": logging.warning,
            "error": logging.error,
            "critical": logging.critical,
        }

        async def wrapper(*args, **kwargs):
            nonlocal levels
            try:
                return await func(*args, **kwargs)
            except catch_exc as e:
                levels[logger_level](f"{func.__name__}, details={e}")
                raise raise_exc(message)

        return wrapper

    return decorator
