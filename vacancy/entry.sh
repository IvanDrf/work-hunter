#!/bin/sh
set -e

echo "Alembic migrations"
uv run alembic upgrade head

echo "Starting application"
uv run -m src.app.main.grpc_server 2>&1 &
uv run -m src.app.main.rabbitmq_consumer 2>&1 &

wait