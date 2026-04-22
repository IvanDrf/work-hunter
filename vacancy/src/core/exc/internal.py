class InternalError(Exception):
    def __init__(self, *args: object) -> None:
        if args:
            self.message = args[0]
        else:
            self.message = None

    def __str__(self) -> str:
        if self.message is not None:
            return f'InternalError: {self.message}'

        return 'InternalError'
