import logging
from asyncio import run

from src.app.app import App
from src.core.config import Config


async def main() -> None:
    config = Config()
    app = App(config)

    await app.init()

    try:
        await app.run()
    finally:
        logging.info("Stopping applications service")
        await app.stop()


if __name__ == "__main__":
    try:
        run(main())
    except KeyboardInterrupt:
        pass
