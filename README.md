## Установка:
1) Для работы, склонируйте репозиторий:
    - git clone git@github.com:Aldirane/go_streaming_app.git
2) Создайте базу данных postgres, подставьте данные для подключения в .env
    - импортируйте таблицы Order, Payment, Delivery, Item. Подставьте ваши данные и путь к файлу:
    -  psql -h localhost -U user -W -d database_name -f go_streaming_app/pkg/database/models.sql
2) Запустите nats streaming server:
    - go run github.com/nats-io/nats-streaming-server
3) Запустите stan publisher для публикаций:
    - go run cmd/stanpub/main.go
4) Запустите сервер и подписчика:
    - go run cmd/server/main.go # Ctrl+c остановит подписку, вторая команда ctrl+c остановит сервер

## API:
    - host:port/ # Получить все ордера из кеша
    - host:port/order_id?order="order id"  # получить определенный ордер по его id

## Структура:
```
.
├── ...
├── workdir/
│     ├── cmd/
│          ├── server/
│                ├── main.go  # сервер и подписчик stan subscribe
│          ├── stanpub/
│                ├──main.go  # stan publish, публикатор
│     ├── pkg/
│          ├── cache/
│                ├── cache.go  # кэш
│          ├── database/
│                ├── postgres/ # Insert, Select запросы
│                ├── database.go  # подключение к БД
│                ├── model.json # json модель ордера
│                ├── models.sql # таблицы для модели ордера
│          ├── order/
│                ├── order.go # структуры Order, Delivery, Payment, Item
│          ├── stansub/
│                ├── stansub.go # Старт подписки stan subscribe
│          ├── test/
│                ├── helper/
│                     ├── testCache.go # тест пакета workdir/pkg/cache
│                     ├── testDatabase.go # тест пакета workdir/pkg/database
│                     ├── testOrder.go # тест пакета workdir/pkg/order
│                ├── main.go # запуск тестов пакета workdir/pkg/test/helper
└── ...
```