import logging
from functools import wraps


logger_levels = {
    'debug': logging.debug,
    'info': logging.info,
    'warning': logging.warning,
    'error': logging.error,
    'critical': logging.critical
}


def catch_rise_error(
    expect_error: tuple[type[Exception]] | type[Exception],
    raise_error: type[Exception],
    logger_level: str, message: str,
):
    def decorator(func):
        @wraps(func)
        async def wrapper(*args, **kwargs):
            try:
                res = await func(*args, **kwargs)
                return res
            except expect_error as e:
                logger_levels[logger_level](f'{func.__name__}: {e}')

                raise raise_error(message)
        return wrapper

    return decorator
