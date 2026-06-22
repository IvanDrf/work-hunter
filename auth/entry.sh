#!/bin/sh
set -e

echo "PostgreSQL migrations"
./migrator \
    --mig=. \
    --cmd=up \
    --steps=2

echo "Starting email worker"
./auth_service &
./email_worker &

wait