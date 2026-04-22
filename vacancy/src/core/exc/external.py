from src.core.exc.base import BaseError


class ExternalError(BaseError):
    def __str__(self) -> str:
        if self.message is not None:
            return f'ExternalError: {self.message}'

        return 'ExternalError'
