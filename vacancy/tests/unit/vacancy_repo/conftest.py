from alembic.command import downgrade as alembic_downgrade
from alembic.command import upgrade as alembic_upgrade
from alembic.config import Config
from pytest import fixture

from src.core.config import PostgreSQLConfig


@fixture(scope="package")
def test_db_config() -> PostgreSQLConfig:
    config = PostgreSQLConfig()

    config.postgres_host = "localhost"
    config.postgres_port = 5432
    config.postgres_user = "test_user"
    config.postgres_password = "test_password"

    config.postgres_db_name = "test_vacancy"

    return config


@fixture(scope="package", autouse=True)
def apply_migrations(test_db_config: PostgreSQLConfig):
    alembic_config = Config("alembic.ini")
    alembic_config.set_main_option("sqlalchemy.url", test_db_config.postgres_dsn)
    alembic_config.set_main_option("script_location", "alembic")

    alembic_upgrade(alembic_config, "head")

    yield

    alembic_downgrade(alembic_config, "base")
