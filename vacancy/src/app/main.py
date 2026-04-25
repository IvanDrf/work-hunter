from asyncio import run
from argparse import ArgumentParser

from src.app.app import App
from src.core.config.config import Config


async def main() -> None:
    parser = ArgumentParser()
    parser.add_argument('--config')
    args, _ = parser.parse_known_args()

    config = Config.load_from_yaml(args.config)

    app = App(config)
    await app.init()
    await app.run()


if __name__ == '__main__':
    run(main())
