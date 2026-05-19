class BaseError(Exception):
    pass


class AccessError(BaseError):
    pass


class ArgumentError(BaseError):
    pass


class InternalError(BaseError):
    pass


class NotFoundError(BaseError):
    pass
