### Order Service
####  Order Service — это проект для управления заказами.
**Технологии**
- **Go**
- **Gorm/PostgreSQL**
- **Kafka**
- **Goose**
- **Docker Compose**

**Установка и запуск**
- Клонируйте репозиторий:
```
git clone https://github.com/Mrsananabos/orderService.git
cd order-service
```
- Запустите сервисы с помощью Docker Compose:
```
docker compose up -d --build   
```
### Использование
**Получение информации о заказе**
```
GET /order/${order_uid}
```
**Добавление заказа**<br>
Чтобы добавить заказ в базу данных, необходимо отправить сообщение в топик Orders. Вы можете сделать это с помощью утилиты kafka-console-producer.sh:
```
kafka-console-producer.sh --bootstrap-server localhost:9092 --topic Orders
```
Сообщение должно содержать информацию о заказе в формате JSON:
```
{
  "order_uid": "4e9ad8fb-2611-46f9-9458-20b59253086b",
  "track_number": "WBILMTESTTRACK",
  "entry": "WBIL",
  "delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
  },
  "payment": {
    "transaction": "b563feb7b2b84b6test",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 1817,
    "payment_dt": 1637907727,
    "bank": "alpha",
    "delivery_cost": 1500,
    "goods_total": 317,
    "custom_fee": 0
  },
  "items": [
    {
      "chrt_id": 9934930,
      "track_number": "WBILMTESTTRACK",
      "price": 453,
      "rid": "ab4219087a764ae0btest",
      "name": "Mascaras",
      "sale": 30,
      "size": "0",
      "total_price": 317,
      "nm_id": 2389212,
      "brand": "Vivienne Sabo",
      "status": 202
    }
  ],
  "locale": "en",
  "internal_signature": "",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}
```