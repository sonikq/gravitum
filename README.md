# Gravitum

Этот репозиторий содержит сервис управления пользователями, разработанный на Go.

## Обзор проекта

Сервис предоставляет функциональность управления пользователями через RESTful API интерфейс.
Он разработан для запуска как непосредственно на локальной машине, так и в Docker-контейнере.

## Предварительные требования
- Go 1.22 или новее
- Docker и Docker Compose (для развертывания в контейнере)
- PostgreSQL (для локального развертывания)

## Начало работы

### Запуск с использованием Docker

Самый простой способ запустить сервис — использовать Docker:

1. Клонируйте репозиторий:
   ```bash
   git clone <url>
   cd gravitum_test_task
   ```
2. Вставьте свои переменные окружения в Dockerfile, а именно строку подключения к БД(DATABASE_DSN)   

3. Соберите и запустите контейнер:
  - docker build -f ./Dockerfile -t gravitum .
  - docker run -d -p 3000:3000 --name gravitum-user-service --restart=always gravitum

4. Сервис будет доступен по адресу http://localhost:3000


## Запуск с использованием Docker Compose

1. Создайте файл docker-compose.yml:

```yaml
version: '3'
services:
  app:
    build: .
    ports:
      - "3000:3000"
    environment:
      - RUN_ADDRESS=0.0.0.0:3000
      - DATABASE_DSN=postgres://postgres:postgres@db:5432/gravitum
      - DB_POOL_WORKERS=75
      - CTX_TIMEOUT=5000
      - LOG_LEVEL=release
      - SERVICE_NAME=user-management
    depends_on:
      - db
  db:
    image: postgres:14
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=gravitum
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

2. Запустите с помощью Docker Compose: 
  - docker-compose up -d

## Локальный запуск:

1. Клонируйте репозиторий:
```bash
git clone <url>
cd gravitum_test_task
   ```

2. Установите зависимости:
  - go mod download

3. Настройте переменные окружения:
```bash
export RUN_ADDRESS=localhost:3000
export DATABASE_DSN=postgres://user:password@localhost:5432/gravitum
export DB_POOL_WORKERS=75
export CTX_TIMEOUT=5000
export LOG_LEVEL=debug
export SERVICE_NAME=user-management
```

4. Запустите проект: 
   - # На Linux: go run /cmd/user_management/main.go
   - # На Windows: go run .\cmd\user_management\main.go

5. Сервис будет доступен по адресу http://localhost:3000


## Переменные окружения: 

| Переменная      | Описание                                                    | Значение по умолчанию                  | 
|-----------------|-------------------------------------------------------------|----------------------------------------|
| RUN_ADDRESS     | Хост и порт для сервиса                                     | localhost:3003                         | 
| DATABASE_DSN    | Строка подключения к PostgreSQL                             | postgres://user:password@host:port/DB? | 
| DB_POOL_WORKERS | Количество рабочих процессов пула соединений с базой данных | 75                                     |
| CTX_TIMEOUT     | Таймаут контекста в миллисекундах                           | 5000                                   |
| LOG_LEVEL       | Уровень логирования (debug, release)                        | release                                |
| SERVICE_NAME    | Название сервиса                                            | user-management                        |


## Документация API:
  - ``GET /api/users/{id}`` - Получение пользователя по ID
  - ``POST /api/users``  - Создание нового пользователя
  - ``PUT /api/users/{id}`` - Обновление пользователя
  - ``DELETE /api/users/{id}`` - Удаление пользователя

# Для создания и обновления нужно указывать Body, пример:
```json
{
    "username": "test_developer",
    "first_name": "Test",
    "last_name": "Testov",
    "email": "test-testov@mail.ru",
    "gender": "M",
    "age": 42
}
```

## Структура проекта: 

gravitum/
├── bin/                  # Скомпилированные бинарные файлы
├── cmd/                  
│   ├──user_management/
│      ├── main.go        # Точка входа в приложение
├── internal/ 
│   ├── app/              # Основной код приложения
│   ├── handler/          # Обработчики API
│   ├── config/           # Конфигурация
│   ├── models/           # Модели данных
│   ├── repository/       # Слой доступа к базе данных
│   ├── service/          # Бизнес-логика
│   └── server/           # HTTP или gRPC сервер
├── pkg/                  # Экспортируемые компоненты
│   ├── logger/           # Логгер
│   ├── reader/           # Обработчик для чтения любых типов данных
│   ├── retrier/          # Пакет для повторного выполнения любых функций
│   └── validator/        # Пакет для валидации данных
├── go.mod                # Определение Go-модуля
├── go.sum                # Контрольные суммы Go-модуля
├── Dockerfile            # Инструкции для сборки Docker
├── docker-compose.yml    # Инструкции для сборки через docker-compose
└── README.md             # Этот файл
