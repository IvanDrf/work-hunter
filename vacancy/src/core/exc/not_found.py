from src.core.exc.base import BaseError


class NotFoundError(BaseError):
    def __str__(self) -> str:
        if self.message is not None:
            return f'NotFound: {self.message}'

        return 'NotFound'
