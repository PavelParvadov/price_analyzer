## Price Analyzer

Cистема потоковой обработки котировок. Генерирует и принимает цены, передаёт их через Kafka, кэширует последние значения в Redis и отдаёт по HTTP.

### Как это работает
- **price-producer**: генерирует случайные цены для заданных тикеров с интервалом (конфигурируется) и/или принимает ручные цены по HTTP, публикует в Kafka.
- **Kafka (3 брокера)**: для доставки цен.
- **price-cache**: консюмит сообщения из Kafka, сохраняет последнюю цену каждого тикера в Redis, предоставляет gRPC метод `GetLatestPrice`.
- **api-gateway**: HTTP-шлюз, вызывает gRPC `price-cache` и отдаёт JSON по HTTP.
- **Redis**: хранит последние цены.


### Стек
- **Go**
- **Kafka** (segmentio/kafka-go)
- **Redis**
- **gRPC** 
- **HTTP**
- **Docker**
- **Docker Compose**
- **zap**
- **cleanenv**
- **net/http**
### Быстрый старт


1) Создайте `.env` в корне проекта, чтобы запустить Redis с паролем:
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


### Эндпойнты и примеры
- **Опубликовать цену вручную (producer)**
  - POST `http://localhost:8080/v1/produce`
  - Body:
    ```json
    {
    "symbol": "AAPL",
     "value": 210.5
    }
    ```


- **Получить последнюю цену (api-gateway)**
  - GET `http://localhost:8081/price?symbol=AAPL`
  
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
- Kafka брокеры: `9091`, `9092`, `9093`
- Kafka UI: `9020`
- Zookeeper: `2181`

### Конфигурация
Можно менять через `config.yaml` в сервисах или переменными окружения


