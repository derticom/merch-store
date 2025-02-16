# merch-store

## Описание

Сервис для управления магазином мерча, где сотрудники могут покупать товары за монеты и делиться монетами друг с другом.

## Регистрация и авторизация

В рамках реализации API, эндпоинт `/api/auth` разделен на два отдельных эндпоинта:
- **Регистрация**: `POST /api/auth/register` — создание нового пользователя.
- **Авторизация**: `POST /api/auth/login` — аутентификация существующего пользователя и получение JWT-токена.

Оба эндпоинта возвращают JWT-токен, который используется для доступа к защищенным ресурсам API.

## Запуск

1. Перед запуском создайте файл `.env` в корневой директории проекта и добавьте в него JWT_SECRET:
   ```JWT_SECRET=super-secure-key```

2. Убедитесь в наличии установленных Docker и Docker Compose.
3. Соберите и запустите контейнеры:
   ```docker-compose up --build```

4. Сервер будет доступен по адресу http://localhost:8080.

## Примеры запросов
#### Регистрация пользователя
   ```
   curl -X POST http://localhost:8080/api/auth/register \
     -H "Content-Type: application/json" \
     -d '{"username": "user1", "password": "password123"}'
   ```

#### Авторизация пользователя
   ```
   curl -X POST http://localhost:8080/api/auth/login \
     -H "Content-Type: application/json" \
     -d '{"username": "user1", "password": "password123"}'
  ```

#### Получение информации о пользователе
   ```
   curl -X GET http://localhost:8080/api/info \
     -H "Authorization: Bearer <your-jwt-token>"
   ```

#### Передача монет
   ```
   curl -X POST http://localhost:8080/api/sendCoin \
     -H "Authorization: Bearer <your-jwt-token>" \
     -H "Content-Type: application/json" \
     -d '{"toUser": "<recipient-user-id>", "amount": 100}'
   ```

#### Покупка товара
   ```
   curl -X GET http://localhost:8080/api/buy/t-shirt \
     -H "Authorization: Bearer <your-jwt-token>"
   ```

## Changelog

### v1.0.0 (16.02.25)
- Сервис создан.