import logging
from functools import wraps

from grpc import StatusCode
from grpc.aio import ServicerContext

from src.core.exc import AccessError, AlreadyExistsError, ArgumentError, InternalError


def handle_errors(func):
    log = logging.getLogger("handle_errors")

    @wraps(func)
    async def wrapper(*args, **kwargs):
        if "context" in kwargs:
            context: ServicerContext = kwargs["context"]

        elif len(args) >= 3:
            context: ServicerContext = args[-1]
        else:
            raise RuntimeError("invalid function in wrapper, msut be like func(..., context: ServicerContext)")

        try:
            return await func(*args, **kwargs)
        except TimeoutError as e:
            log.critical(f"TimeoutError {func.__name__}, details={e}")
            await context.abort(code=StatusCode.INTERNAL, details="timeout, something is too long")

        except ArgumentError as e:
            log.info(f"ArgumentError {func.__name__}, details={e}")
            await context.abort(code=StatusCode.INVALID_ARGUMENT, details=str(e))

        except InternalError as e:
            log.critical(f"InternalError {func.__name__}, details={e}")
            await context.abort(code=StatusCode.INTERNAL, details=str(e))

        except AccessError as e:
            log.critical(f"AccessError {func.__name__}, details={e}")
            await context.abort(code=StatusCode.PERMISSION_DENIED, details=str(e))

        except AlreadyExistsError as e:
            log.info(f"AlreadyExistsError {func.__name__}, details={e}")
            await context.abort(code=StatusCode.ALREADY_EXISTS, details=str(e))

    return wrapper


MIN_LIMIT = 1
MAX_LIMIT = 50

MIN_OFFSET = 0


def validate_limit_offset(func):
    @wraps(func)
    async def wrapper(*args, **kwargs):
        request = args[1] or kwargs["request"]

        if not hasattr(request, "limit"):
            raise RuntimeError(f"request doesn't has limit in body, {request=}")

        if not hasattr(request, "offset"):
            raise RuntimeError(f"request doesn't has offset in body {request=}")

        if not (MIN_LIMIT <= request.limit <= MAX_LIMIT):
            raise ArgumentError(f"invalid limit value in request, limit={request.limit}, but must be in {MIN_LIMIT}-{MAX_LIMIT}")

        if not request.offset >= MIN_OFFSET:
            raise ArgumentError(f"invalid offset value in request must be greater than {MIN_OFFSET}, but offset={request.offset}")

        return await func(*args, **kwargs)

    return wrapper
