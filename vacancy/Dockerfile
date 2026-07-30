FROM python:3.12-slim-trixie
COPY --from=ghcr.io/astral-sh/uv:latest /uv /uvx /bin/

ENV PYTHONDONTWRITEBYTECOD=1
ENV PYTHONUNBUFFERED=1
ENV UV_NO_DEV=1

WORKDIR /app

COPY vacancy/.python-version .
COPY vacancy/pyproject.toml .
COPY vacancy/uv.lock .
RUN uv sync --locked

COPY pkg/vacancy_api pkg/vacancy_api
COPY pkg/common pkg/common

COPY vacancy/ .
RUN chmod +x entry.sh

EXPOSE 50052

ENTRYPOINT ["./entry.sh"]
