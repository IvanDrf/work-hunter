import logging

levels = {
    "debug": logging.DEBUG,
    "info": logging.INFO,
    "warning": logging.WARNING,
    "error": logging.ERROR,
    "critical": logging.CRITICAL,
}


def setup_logger(logger_level: str, log_file: str | None = None) -> None:
    logging.basicConfig(
        level=levels[logger_level.lower()], filename=log_file, filemode="a", format="%(asctime)s %(levelname)s %(message)s"
    )
