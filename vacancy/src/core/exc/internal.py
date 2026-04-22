from src.core.exc.base import BaseError


class InternalError(BaseError):
    def __str__(self) -> str:
        if self.message is not None:
            return f'InternalError: {self.message}'

        return 'InternalError'
