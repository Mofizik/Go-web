# Go gRPC Server
Этот проект создан с целью научиться веб разработке продвинутых систем на Go.
Сервер используется для взаимодействия с заказами через gRPC протокол.

# Содержание
1. [Подготовка](#подготовка)
2. [Запуск](#запуск)
3. [Документация](#документация)
    1. [Структура](#структура)
    2. [Методы API](#методы-api)

## Подготовка

1. Установите зависимости:
    - [Go 1.21+](https://golang.org/dl/)
    - [Protocol Buffers (protoc)](https://grpc.io/docs/protoc-installation/)
    - [grpcurl](https://github.com/fullstorydev/grpcurl) — для тестирования API

2. Создайте `.env` файл в `./internal/order/config/` на основе примера:
```bash
cp internal/order/config/.env.example internal/order/config/.env
```

Переменные окружения:

| Переменная | Описание                                  | Пример |
|------------|-------------------------------------------|--------|
| GRPC_PORT  | Порт, на котором будет запущен сервер     | 50051  |
| APP_ENV    | Среда разработки (`local`, `dev`, `prod`) | local  |

## Запуск

Для сборки и запуска сервера используйте Makefile:

```bash
make        # генерация proto, сборка и запуск
```

Доступные команды:

```bash
make protoc      # генерация gRPC кода из order.proto
make build      # сборка бинарника
make run        # запуск сервера
```

После первой сборки сервер можно запускать напрямую из корня проекта:

```bash
./main
```

## Документация

### Структура

```
order-service/
├── cmd/server/main.go          # точка входа
├── internal/
│   ├── app/app.go              # инициализация и запуск приложения
│   └── order/
│       ├── api/handler/        # gRPC хендлеры
│       ├── config/             # .env и .env.example
│       ├── model/              # модель Order
│       ├── service/            # бизнес-логика
│       └── storage/            # хранилище данных (sync.Mutex)
├── pkg/
│   ├── api/test/               # сгенерированный protobuf код
│   ├── config/                 # загрузка конфигурации
│   ├── idgen/                  # генерация уникальных ID
│   └── logger/                 # настройка логгера
├── order.proto                 # описание gRPC сервиса
├── makefile
├── go.mod
└── go.sum
```

### Методы API

Для тестирования используется `grpcurl`. Убедитесь что сервер запущен.

---

#### CreateOrder — создание заказа

Принимает:

| Поле       | Тип    | Обязательное | Описание        |
|------------|--------|:------------:|-----------------|
| `item`     | string | да           | Название товара |
| `quantity` | int32  | да           | Количество      |

Возвращает:

| Поле | Тип    | Описание                        |
|------|--------|---------------------------------|
| `id` | string | Уникальный ID созданного заказа |

```bash
grpcurl -plaintext -d '{"item": "apple", "quantity": 5}' \
  localhost:50051 api.OrderService/CreateOrder
```

```json
{
  "id": "a1b2c3d4e5f6g7h8"
}
```

---

#### GetOrder — получить заказ по ID

Принимает:

| Поле | Тип    | Обязательное | Описание  |
|------|--------|:------------:|-----------|
| `id` | string | да           | ID заказа |

Возвращает:

| Поле             | Тип    | Описание        |
|------------------|--------|-----------------|
| `order.id`       | string | ID заказа       |
| `order.item`     | string | Название товара |
| `order.quantity` | int32  | Количество      |

```bash
grpcurl -plaintext -d '{"id": "a1b2c3d4e5f6g7h8"}' \
  localhost:50051 api.OrderService/GetOrder
```

```json
{
  "order": {
    "id": "a1b2c3d4e5f6g7h8",
    "item": "apple",
    "quantity": 5
  }
}
```

---

#### UpdateOrder — обновить заказ

Принимает:

| Поле       | Тип    | Обязательное | Описание                         |
|------------|--------|:------------:|----------------------------------|
| `id`       | string | да           | ID заказа который нужно изменить |
| `item`     | string | да           | Новое название товара            |
| `quantity` | int32  | да           | Новое количество                 |

Возвращает обновлённый заказ:

| Поле             | Тип    | Описание        |
|------------------|--------|-----------------|
| `order.id`       | string | ID заказа       |
| `order.item`     | string | Название товара |
| `order.quantity` | int32  | Количество      |

```bash
grpcurl -plaintext -d '{"id": "a1b2c3d4e5f6g7h8", "item": "banana", "quantity": 10}' \
  localhost:50051 api.OrderService/UpdateOrder
```

```json
{
  "order": {
    "id": "a1b2c3d4e5f6g7h8",
    "item": "banana",
    "quantity": 10
  }
}
```

---

#### DeleteOrder — удалить заказ

Принимает:

| Поле | Тип    | Обязательное | Описание  |
|------|--------|:------------:|-----------|
| `id` | string | да           | ID заказа |

Возвращает:

| Поле      | Тип  | Описание                                    |
|-----------|------|---------------------------------------------|
| `success` | bool | `true` если заказ удалён, `false` если нет |

```bash
grpcurl -plaintext -d '{"id": "a1b2c3d4e5f6g7h8"}' \
  localhost:50051 api.OrderService/DeleteOrder
```

```json
{
  "success": true
}
```

---

#### ListOrders — список всех заказов

Принимает: пустой запрос `{}`

Возвращает:

| Поле     | Тип     | Описание            |
|----------|---------|---------------------|
| `orders` | []Order | Массив всех заказов |

Каждый элемент массива `Order`:

| Поле       | Тип    | Описание        |
|------------|--------|-----------------|
| `id`       | string | ID заказа       |
| `item`     | string | Название товара |
| `quantity` | int32  | Количество      |

```bash
grpcurl -plaintext -d '{}' \
  localhost:50051 api.OrderService/ListOrders
```

```json
{
  "orders": [
    {
      "id": "a1b2c3d4e5f6g7h8",
      "item": "banana",
      "quantity": 10
    }
  ]
}
```