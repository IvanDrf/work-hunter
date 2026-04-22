from logging import CRITICAL, DEBUG, ERROR, INFO, WARNING, basicConfig


levels = {
    'debug': DEBUG,
    'info': INFO,
    'warn': WARNING,
    'error': ERROR,
    'critical': CRITICAL,
}


def setup_logger(level: str, log_file: str | None = None) -> None:
    if log_file:
        basicConfig(
            level=levels[level], filename=log_file, filemode="a",
            format="%(asctime)s %(levelname)s %(message)s"
        )

    else:
        basicConfig(
            level=levels[level], format="%(asctime)s %(levelname)s %(message)s"
        )
