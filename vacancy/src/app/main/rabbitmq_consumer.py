from asyncio import run

from src.app.app import ConsumerApp
from src.core.config import Config
from src.core.logger.setup import setup_logger


async def main() -> None:
    config = Config()
    setup_logger(config.logger_level)

    app = ConsumerApp(config)
    await app.init()
    try:
        await app.run()
    finally:
        await app.stop()


if __name__ == "__main__":
    try:
        run(main())
    except KeyboardInterrupt:
        pass
