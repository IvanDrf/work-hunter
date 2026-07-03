#!/bin/sh
set -e

cat .env.example > .env
cd tests && docker compose up -d && cd ..
uv run -m pytest tests/ -v -ss
cd tests && docker compose down -v