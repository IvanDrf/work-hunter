class MessageBox[T]:
    def __init__(self, max_size: int) -> None:
        self.box: list[T | None] = [None] * max_size
        self.pos: int = 0

    def add_message(self, message: T) -> None:
        if self.pos > len(self.box) - 1:
            raise OverflowError(f"box is already full, max_size={len(self.box)}, pos={self.pos + 1}")

        self.box[self.pos] = message
        self.pos += 1

    def drop_box(self) -> None:
        self.box = [None] * self.pos
        self.pos = 0

    def get_messages(self) -> list[T | None]:
        return self.box
