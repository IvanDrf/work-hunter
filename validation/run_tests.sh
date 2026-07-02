#!/bin/sh
set -e

cd tests && docker-compose up -d && cd ..
uv run -m pytest tests/ -v -ss
cd tests && docker-compose down -v