#!/bin/sh
set -e

echo "Alembic migrations"
uv run alembic upgrade head

echo "Starting application"
uv run -m src.app.main
