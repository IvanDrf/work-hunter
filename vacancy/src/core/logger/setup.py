from logging import CRITICAL, DEBUG, ERROR, INFO, WARNING, basicConfig

levels = {
    "debug": DEBUG,
    "info": INFO,
    "warn": WARNING,
    "error": ERROR,
    "critical": CRITICAL,
}


def setup_logger(level: str, log_file: str | None = None) -> None:
    level = level.lower()

    basicConfig(level=levels[level], filename=log_file, format="%(asctime)s %(name)s %(levelname)s %(message)s")
