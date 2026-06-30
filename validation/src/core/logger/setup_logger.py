import logging

levels = {
    "debug": logging.DEBUG,
    "info": logging.INFO,
    "warn": logging.WARNING,
    "error": logging.ERROR,
    "critical": logging.CRITICAL,
}


def setup_logger(level: str) -> None:
    level = level.lower()

    logging.basicConfig(level=levels[level], format="%(asctime)s %(levelname)s %(message)s")
