#!/bin/sh
set -e

echo "Alembic migrations"
alembic upgrade head

echo "Starting application"
exec python3 -m src.app.main