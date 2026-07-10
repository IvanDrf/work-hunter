FROM python:3.12-slim-trixie
COPY --from=ghcr.io/astral-sh/uv:latest /uv /uvx /bin/

ENV PYTHONDONTWRITEBYTECOD=1
ENV PYTHONUNBEFFERED=1
ENV UV_NO_DEV=1

WORKDIR /app

COPY applications/.python-version .
COPY applications/pyproject.toml .
COPY applications/uv.lock .
RUN uv sync --locked

COPY pkg/applications_api pkg/applications_api
COPY pkg/common pkg/common

COPY applications/ .
RUN chmod +x entry.sh

EXPOSE 50053

ENTRYPOINT ["./entry.sh"]
