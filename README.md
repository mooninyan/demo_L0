# demoL0

Простая демонстрационная сервисная служба на Go, построенная по принципам Clean Architecture. Демонстрирует интеграцию с PostgreSQL и Apache Kafka, содержит REST API и потребителей Kafka.

## Основное описание

Сервис предоставляет CRUD-операции для сущности "Order" через REST API, использует GORM для работы с PostgreSQL и интегрируется с Kafka для отправки и получения сообщений. В проекте также реализован кэш в памяти и обработка DLQ (dead-letter queue).

## Архитектура

Проект построен по принципам Clean Architecture:

```
cmd/server/           # Точка входа приложения
internal/
├── domain/           # Бизнес-сущности (MainOrder, Delivery, Payment, Item)
├── ports/            # Интерфейсы (repository, kafka)
├── usecase/          # Бизнес-логика (order service)
├── adapters/         # Внешние зависимости
│   ├── repository/   # PostgreSQL репозитории
│   ├── kafka/        # Kafka producer/consumer
│   ├── migrations/   # SQL миграции
│   ├── cache/        # In-memory кэш
│   └── listener/     # Kafka listeners
├── transport/        # HTTP handlers
├── front/            # Статические страницы
└── utils/            # Вспомогательные функции
```

## Особенности

- REST API для управления заказами
- Хранение данных в PostgreSQL через GORM
- Интеграция с Kafka (producer/consumer, создание топиков)
- Поддержка DLQ (опционально через переменную окружения)
- Миграции базы данных в `pkg/database/migrations`
- Простая страница фронтенда в `internal/front`

## Требования

- Go 1.20+ (проверьте `go.mod`)
- PostgreSQL
- Apache Kafka

## Переменные окружения

Проект загружает переменные из `.env` (используется `github.com/joho/godotenv`). Основные переменные, которые используются в коде:

- `DB_USER` — пользователь PostgreSQL
- `DB_PASSWORD` — пароль PostgreSQL
- `DB_NAME` — имя базы данных
- `DB_PORT` — порт PostgreSQL
- `KAFKA_HOST_PORT` — адрес Kafka-брокера (например, `localhost:9092`)
- `KAFKA_TOPIC` — имя основного Kafka-топика
- `KAFKA_DLQ_TOPIC` — имя DLQ-топика
- `KAFKA_TOPIC_PARTITION_COUNT` — количество партиций для основного топика
- `KAFKA_DLQ_TOPIC_PARTITION_COUNT` — количество партиций для DLQ
- `KAFKA_GROUP_ID` — consumer group id для основного потребителя
- `KAFKA_DLQ_GROUP_ID` — consumer group id для DLQ
- `LISTEN_DLQ` — `true`/`false`, слушать DLQ или нет

## Установка

1. Клонируйте репозиторий

```bash
git clone <repo-url>
cd demoL0
```

2. Установите зависимости

```bash
go mod download
```

3. Настройте `.env` файл с переменными окружения (см. раздел выше)

4. Запустите PostgreSQL и Kafka, убедитесь, что указанные в `.env` адреса корректны

## Запуск

Простой способ запустить сервис в development режиме:

```bash
go run cmd/server/main.go
```

Сервис по умолчанию слушает HTTP на порту `8081`.

### Переменные окружения

Скопируйте `docs/env.example` в `.env` и настройте под вашу среду:

```bash
cp docs/env.example .env
```

## Миграции

Файлы миграций находятся в `internal/adapters/migrations`. При старте приложение вызывает `InitialMigration()` и `RunMigrations()` — убедитесь, что у вас есть права на создание/изменение схемы в указанной базе.

## API

REST-ендпойнты реализованы в `internal/transport/http/handler.go`. Основные операции (CRUD) для заказов доступны через HTTP. Используется роутер Gorilla Mux.

### Доступные эндпойнты:

- `GET /order` - получить все заказы
- `GET /order/{id}` - получить заказ по ID
- `POST /order` - создать новый заказ
- `PUT /order/{id}` - обновить заказ
- `DELETE /order/{id}` - удалить заказ

## Kafka

- При запуске вызывается `InitializeTopics()` — создаются (или обновляются) топики, указанные в `.env`.
- Создаётся producer (`CreateKafkaWriter`) и consumer group (`CreateConsumerGroup`, `CreateDlqConsumerGroup`).

## Как тестировать

- Напишите и выполните unit-тесты с помощью `go test ./...`
- Для интеграционных тестов используйте развёрнутые инстансы PostgreSQL и Kafka

## Разработка и вклад

PR приветствуются. Откройте issue для обсуждения больших изменений.

