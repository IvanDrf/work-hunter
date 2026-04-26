import logging
from functools import wraps

from grpc import ServicerContext, StatusCode

from src.core.exc import AccessError, ArgumentError, InternalError, NotFoundError


def handle_errors(func):
    @wraps(func)
    async def wrapper(*args, **kwargs):
        context: ServicerContext = args[2] or kwargs['context']
        try:
            res = await func(*args, **kwargs)
            return res

        except ArgumentError as e:
            logging.info(f'{func.__name__}: {e}')

            await context.abort(StatusCode.INVALID_ARGUMENT, e.__str__())

        except AccessError as e:
            logging.info(f'{func.__name__}: {e}')

            await context.abort(StatusCode.PERMISSION_DENIED, e.__str__())

        except InternalError as e:
            logging.critical(f'{func.__name__}: {e}')

            await context.abort(StatusCode.INTERNAL, e.__str__())

        except NotFoundError as e:
            logging.info(f'{func.__name__}: {e}')

            await context.abort(StatusCode.NOT_FOUND, e.__str__())

    return wrapper
