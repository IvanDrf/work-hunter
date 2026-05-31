FROM python:3.12-slim

ENV PYTHONDONTWRITEBYTECOD=1
ENV PYTHONUNBEFFERED=1

WORKDIR /app

COPY vacancy/requirements.txt . 
RUN pip install --no-cache-dir -r requirements.txt

COPY pkg/vacancy_api pkg/vacancy_api
COPY pkg/common pkg/common


COPY vacancy/ .
RUN chmod +x entry.sh

EXPOSE 50052

ENTRYPOINT ["./entry.sh"]
