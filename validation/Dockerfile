FROM python:3.12-slim-trixie
COPY --from=ghcr.io/astral-sh/uv:latest /uv /uvx /bin/

ENV PYTHONDONTWRITEBYTECOD=1
ENV PYTHONUNBEFFERED=1
ENV UV_NO_DEV=1

WORKDIR /app

COPY .python-version .
COPY pyproject.toml .
COPY uv.lock .
RUN uv sync --locked


COPY . .

EXPOSE 8080

CMD ["uv", "run", "-m", "src.app.main"]
