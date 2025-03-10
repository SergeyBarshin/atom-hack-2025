FROM python:3.11-slim

WORKDIR /app

COPY ./analytics-service/requirements.txt /app/requirements.txt
RUN pip install --no-cache-dir -r requirements.txt

# Копируем остальные файлы
COPY ./analytics-service /app


EXPOSE 5000

ENV FLASK_APP=app.py

CMD ["flask", "run", "--host=0.0.0.0", "--port=5000"]
