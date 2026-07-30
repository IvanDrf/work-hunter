#!/bin/sh
set -e

echo "=== Running PostgreSQL migrations ==="
./migrator \
    -mig=/app/migrations \
    -cmd=up \
    -steps=1

echo "=== Starting application server ==="
exec ./server