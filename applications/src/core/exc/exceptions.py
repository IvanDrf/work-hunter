class BaseError(Exception):
    pass


class ArgumentError(BaseError):
    pass


class InternalError(BaseError):
    pass


class AccessError(BaseError):
    pass


class AlreadyExistsError(BaseError):
    pass


class NotFoundError(BaseError):
    pass
