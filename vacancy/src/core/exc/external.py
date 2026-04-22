class ExternalError(Exception):
    def __init__(self, *args: object) -> None:
        if args:
            self.message = args[0]
        else:
            self.message = None

    def __str__(self) -> str:
        if self.message is not None:
            return f'ExternalError: {self.message}'

        return 'ExternalError'
