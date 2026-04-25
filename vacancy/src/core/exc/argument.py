from src.core.exc.base import BaseError


class ArgumentError(BaseError):
    def __str__(self) -> str:
        if self.message is not None:
            return f'ArgumentError: {self.message}'

        return 'ArgumentError'
