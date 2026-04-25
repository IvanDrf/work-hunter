import logging
from functools import wraps

from grpc import ServicerContext, StatusCode

from src.core.exc.access import AccessError
from src.core.exc.argument import ArgumentError
from src.core.exc.internal import InternalError
from src.core.exc.not_found import NotFoundError


def handle_errors(func):
    @wraps(func)
    async def wrapper(*args, **kwargs):
        context: ServicerContext = args[1] or kwargs['context']
        try:
            res = await func(*args, **kwargs)
            return res

        except ArgumentError as e:
            logging.info(f'{func.__name__}: {e}')

            context.abort(StatusCode.INVALID_ARGUMENT, e.__str__())

        except AccessError as e:
            logging.info(f'{func.__name__}: {e}')

            context.abort(StatusCode.PERMISSION_DENIED, e.__str__())

        except InternalError as e:
            logging.critical(f'{func.__name__}: {e}')

            context.abort(StatusCode.INTERNAL, e.__str__())

        except NotFoundError as e:
            logging.info(f'{func.__name__}: {e}')

            context.abort(StatusCode.NOT_FOUND, e.__str__())

    return wrapper
