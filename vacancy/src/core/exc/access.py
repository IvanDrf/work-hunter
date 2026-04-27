from src.core.exc.base import BaseError


class AccessError(BaseError):
    def __str__(self) -> str:
        if self.message:
            return f'AccessError: {self.message}'

        return 'AccessError'
