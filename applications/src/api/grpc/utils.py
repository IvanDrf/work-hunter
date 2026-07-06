import logging
from functools import wraps

from grpc import StatusCode
from grpc.aio import ServicerContext
from src.core.exc import AccessError, ArgumentError, InternalError


def handle_errors(func):
    @wraps(func)
    async def wrapper(*args, **kwargs):
        if "context" in kwargs:
            context: ServicerContext = kwargs["context"]

        elif len(args) >= 2 and isinstance(args[-1], ServicerContext):
            context: ServicerContext = args[-1]
        else:
            raise RuntimeError("invalid function in wrapper, msut be like func(..., context: ServicerContext)")

        try:
            return await func(*args, **kwargs)
        except TimeoutError as e:
            logging.critical(f"TimeoutError {func.__name__}, details={e}")
            await context.abort(code=StatusCode.INTERNAL, details="timeout, something is too long")

        except ArgumentError as e:
            logging.info(f"ArgumentError {func.__name__}, details={e}")
            await context.abort(code=StatusCode.INVALID_ARGUMENT, details=str(e))

        except InternalError as e:
            logging.critical(f"InternalError {func.__name__}, details={e}")
            await context.abort(code=StatusCode.INTERNAL, details=str(e))

        except AccessError as e:
            logging.critical(f"AccessError {func.__name__}, details={e}")
            await context.abort(code=StatusCode.PERMISSION_DENIED, details=str(e))

    return wrapper
