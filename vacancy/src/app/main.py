from asyncio import run

from src.app.app import App
from src.core.config import Config
from src.core.logger.setup import setup_logger


async def main() -> None:
    config = Config()
    setup_logger(config.logger_level)

    app = App(config)
    await app.init()
    await app.run()


if __name__ == "__main__":
    run(main())
