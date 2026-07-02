#!/bin/sh
set -e

echo API_KEY="TEST_KEY" > .env
cd tests && docker compose up -d && cd ..
uv run -m pytest tests/ -v -ss
cd tests && docker compose down -v