# ethno-TulaHack-2026

Онлайн платформа для популяризации промыслов народов России

# Запуск из Docker контейнера

1. Из корня репозитория запустите все сервисы в фоновом режиме:

```bash
docker compose up -d --build
```

2.  Ждём 5–15 сек, пока PostgreSQL завершит инициализацию

3.  Создаём дамп для корректного изображения интерактивной карты и квестов

Linux/macOS

```bash
docker compose exec -T db pg_dump -U dev -d ethno --format=plain --encoding=UTF8 --no-owner --no-acl --clean --if-exists > seed.sql
```

Windows

```bash
docker compose exec -T db pg_dump -U dev -d ethno --format=plain --encoding=UTF8 --no-owner --no-acl --clean --if-exists | Set-Content -Encoding UTF8 seed.sql
```
