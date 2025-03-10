# Используем минимальный образ Python
FROM python:3.11-slim

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта
COPY . .

# Устанавливаем зависимости
RUN pip install --no-cache-dir -r requirements.txt

# Открываем порт 5000
EXPOSE 5000

# Устанавливаем переменную окружения FLASK_APP
ENV FLASK_APP=app.py

# Запускаем приложение
CMD ["flask", "run", "--host=0.0.0.0", "--port=5000"]
