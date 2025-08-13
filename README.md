## Price Analyzer

Небольшая учебная система потоковой обработки котировок. Генерирует или принимает цены, передаёт их через Kafka, кэширует последние значения в Redis и отдаёт по HTTP.

### Как это работает
- **price-producer**: генерирует случайные цены для заданных тикеров с интервалом и/или принимает ручные цены по HTTP, публикует в Kafka.
- **Kafka (3 брокера + Zookeeper)**: шина сообщений для доставки цен.
- **price-cache**: консюмит сообщения из Kafka, сохраняет последнюю цену каждого тикера в Redis, предоставляет gRPC метод `GetLatestPrice`.
- **api-gateway**: HTTP-шлюз, вызывает gRPC `price-cache` и отдаёт JSON по HTTP.
- **Redis**: хранит последние цены.

Поток: Producer → Kafka → price-cache (Redis) → api-gateway → Клиент.

### Стек
- **Go** (все три сервиса)
- **Kafka** (segmentio/kafka-go)
- **Redis**
- **gRPC** (между api-gateway и price-cache)
- **HTTP** (внешние REST-точки)
- **Docker Compose** для локального запуска
- **zap** для логирования

### Быстрый старт (Docker Compose)
Требования: установлен Docker Desktop.

1) Создайте `.env` в корне проекта (или задайте переменную среды), чтобы запустить Redis с паролем:
```
REDIS_PASSWORD=pavel
```

2) Запустите всё:
```
docker compose up -d --build
```

3) Проверьте UI и точки:
- Kafka UI: `http://localhost:9020`
- Producer HTTP: `http://localhost:8080`
- API Gateway: `http://localhost:8081`

Остановить и очистить:
```
docker compose down -v
```

### Эндпойнты и примеры
- **Опубликовать цену вручную (producer)**
  - POST `http://localhost:8080/v1/produce`
  - Body:
    ```json
    { "symbol": "AAPL", "value": 210.5 }
    ```
  - Пример:
    ```bash
    curl -X POST http://localhost:8080/v1/produce \
      -H "Content-Type: application/json" \
      -d '{"symbol":"AAPL","value":210.5}'
    ```

- **Получить последнюю цену (api-gateway)**
  - GET `http://localhost:8081/price?symbol=AAPL`
  - Пример:
    ```bash
    curl "http://localhost:8081/price?symbol=AAPL"
    ```
  - Ответ (пример):
    ```json
    {
      "exists": true,
      "price": { "symbol": "AAPL", "value": 210.5, "timestamp": "2024-01-01T12:00:00Z" }
    }
    ```

### Порты по умолчанию
- Producer HTTP: `8080`
- API Gateway HTTP: `8081`
- price-cache gRPC: `6000`
- Kafka брокеры (наружу): `9091`, `9092`, `9093`
- Kafka UI: `9020`
- Zookeeper: `2181`

### Конфигурация
Можно менять через `config.yaml` в сервисах или переменными окружения (Compose уже задаёт ключевые):

- Producer (`price-producer`):
  - `ADDRESSES` — адреса брокеров Kafka (пример: `kafka1:29091,kafka2:29092,kafka3:29093`)
  - `TOPIC` — топик (по умолчанию `price`)
  - `TICKERS` — список тикеров (пример: `AAPL,GOOG,TSLA`)
  - `INTERVAL_MS` — интервал генерации в мс (пример: `5000`)
  - `INITIAL_PRICE` — стартовая цена (пример: `200.0`)
  - `VOLATILITY_PERCENT` — волатильность, % (пример: `0.5`)
  - `PRICE_PORT` — порт HTTP сервера (по умолчанию `8080`)

- Price-cache (`price-cache`):
  - `ADDRESSES`, `TOPIC`, `GROUP_ID` — настройки Kafka
  - `REDIS_ADDR` — адрес Redis (в Compose: `redis:6379`)
  - `REDIS_PASSWORD` — пароль Redis (см. `.env`)
  - `REDIS_DB` — номер БД (по умолчанию `0`)
  - `GRPC_PORT` — порт gRPC (по умолчанию `6000`)

- API Gateway (`api-gateway`):
  - `HTTP_PORT` — порт HTTP (по умолчанию `8081`)
  - `PRICE_CACHE_HOST` — хост gRPC (`price-cache` в сети Compose)
  - `PRICE_CACHE_PORT` — порт gRPC (`6000`)

Файлы конфигурации по умолчанию:
- `price-producer/config/config.yaml`
- `price-cache/config/config.yaml`
- `api-gateway/config/config.yaml`

### Полезно знать
- После запуска подождите пару секунд: сервисы имеют healthchecks.
- Для локальных тестов без ручной публикации Producer сам генерирует цены с заданным интервалом.
- В Kafka UI можно смотреть сообщения в топике `price`.


