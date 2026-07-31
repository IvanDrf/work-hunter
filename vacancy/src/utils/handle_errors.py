import logging
from functools import wraps

from grpc import ServicerContext, StatusCode
from pydantic import ValidationError

from src.core.exc import AccessError, ArgumentError, InternalError, NotFoundError


def handle_errors(func):
    logger = logging.getLogger("handle_errors")

    @wraps(func)
    async def wrapper(*args, **kwargs):
        if "context" in kwargs:
            context: ServicerContext = kwargs["context"]
        elif len(args) >= 3:
            context: ServicerContext = args[2]
        else:
            raise ValueError("invalid function signature, must has 'context'")

        try:
            res = await func(*args, **kwargs)
            return res

        except (ArgumentError, ValidationError) as e:
            logger.info(f"{func.__name__}: {e}")

            await context.abort(StatusCode.INVALID_ARGUMENT, e.__str__())

        except AccessError as e:
            logger.info(f"{func.__name__}: {e}")

            await context.abort(StatusCode.PERMISSION_DENIED, e.__str__())

        except InternalError as e:
            logger.critical(f"{func.__name__}: {e}")

            await context.abort(StatusCode.INTERNAL, e.__str__())

        except NotFoundError as e:
            logger.info(f"{func.__name__}: {e}")

            await context.abort(StatusCode.NOT_FOUND, e.__str__())

        except OSError as e:
            logger.critical(f"{func.__name__}: critical os error {e.__str__()}")

            await context.abort(StatusCode.INTERNAL, "internal error, service in not available now")

    return wrapper
