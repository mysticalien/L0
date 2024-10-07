# L0 Orders Service

---

### 🎥 Видеодемонстрация работы сервиса:



---

## Описание проекта

**L0 Orders Service** — это демонстрационный сервис для получения и отображения данных о заказах. Он использует **Kafka** для обработки сообщений о новых заказах, которые затем сохраняются в базе данных **PostgreSQL** и кэшируются в оперативной памяти. Сервис поддерживает REST API для получения данных о заказах по их `order_id`, а также восстановление кэша после перезапуска.

<p align="center">
  <img src="/images/ui_1.png" alt="Orders Service" width="600" />
</p>

---

## Основные функции

1. Подключение к **Kafka** и обработка сообщений о заказах через топик `orders`.
2. Сохранение заказов в **PostgreSQL**.
3. **Кэширование** данных в памяти для быстрого доступа.
4. Восстановление кэша из базы данных после перезапуска.
5. Запуск **HTTP-сервера** для получения данных по ID заказа через REST API.
6. Поддержка **автотестов** и **стресс-тестов**.

---

## Архитектура

Проект состоит из следующих компонентов:

- **Kafka Consumer**: Подписывается на топик `orders`, обрабатывает и сохраняет сообщения.
- **PostgreSQL**: Хранение данных о заказах.
- **Кэш**: In-memory хранение данных о заказах для быстрого доступа.
- **HTTP-сервер**: Позволяет получать данные о заказе через REST API.
- **Redpanda Console**: UI для мониторинга Kafka.

### Технологии:

- **Go**: Язык программирования для реализации сервиса.
- **Kafka**: Обработка сообщений о заказах.
- **PostgreSQL**: Хранение заказов.
- **Docker**: Контейнеризация всех сервисов.
- **Chi**: HTTP роутинг.
- **WRK** и **Vegeta**: Стресс-тестирование.

---

## Как запустить проект

### 1. Клонирование репозитория

```bash
git clone git@github.com:mysticalien/L0.git
cd L0
```

### 2. Запуск с помощью Makefile

Запуск всех сервисов:

```bash
make all
```

Эта команда:
- Соберет Go-приложение.
- Запустит сервисы через Docker Compose.
- Запустит HTTP-сервер на порту `8080`.
- Запустит Kafka producer для отправки сообщения.

### 3. Остановка всех сервисов

```bash
make docker-down
```

---

## Графический интерфейс
Сервис включает простой интерфейс для отображения данных о заказах. Вы можете открыть http://localhost:8080, чтобы увидеть страницу интерфейса.

Интерфейс включает следующие возможности:
- Главная страница: Показ списка заказов и поиск заказа по ID.
- Поиск заказа: Введите order_uid для поиска и просмотра подробной информации о заказе.
- Детализация заказа: Отображение информации о доставке, товарах и оплате заказа.
Пример скриншота страницы:

<p align="center">
  <img src="/images/ui.png" alt="Пользовательский интерфейс" width="600" />
</p>

<p align="center">
  <img src="/images/ui_2.png" alt="Обработка неверного id заказа" width="600" />
</p>

## Команды Makefile

- `make all`: Полный цикл сборки и запуска сервиса, включая Kafka, PostgreSQL, продюсера и потребителя.
- `make create-topic`: Создание топика Kafka вручную.
- `make build`: Сборка основного приложения и продюсера.
- `make docker-up`: Запуск сервисов через Docker Compose.
- `make run-main`: Запуск основного сервиса.
- `make run-producer`: Запуск продюсера Kafka.
- `make docker-down`: Остановка всех сервисов.
- `make clean`: Очистка скомпилированных бинарных файлов.
- `make redpanda-console`: Запуск Redpanda Console для мониторинга Kafka.
- `make test`: Запуск тестов.
- `make test_wrk`: Стресс тестирование с утилитой wrk.
- `make test_vegeta`: Стресс тестирование с утилитой vegeta.

---

## API

1. **Получение данных о заказе по ID**:

```bash
GET http://localhost:8080/order/{id}
```

Пример запроса:

```bash
curl http://localhost:8080/order/b563feb7b2b84b6test
```

2. **Отправка нового заказа через Kafka**:

Producer отправляет данные в топик `orders` через Kafka:

```bash
make run-producer
```

---

## Автотесты

Автотесты включают проверку кэша, работы с базой данных и Kafka. Для их запуска:

```bash
make test
```

Тесты находятся в папке `tests/`.

<p align="center">
  <img src="/images/test.png" alt="Автоматические тесты" width="600" />
</p>

---

## Стресс-тесты

### WRK:

Пример стресс-теста для HTTP-сервера:

```bash
make test_wrk
```

<p align="center">
  <img src="/images/wrk_test.png" alt="WRK тесты" width="600" />
</p>

### Vegeta:

Запускаем тестирование:

```bash
make test_vegeta
```

<p align="center">
  <img src="/images/vegeta_test.png" alt="Vegeta тесты" width="600" />
</p>

---

## Логирование

Все логи сохраняются в формате JSON и выводятся в консоль. Логирование настроено с помощью пакета `slog` и поддерживает различные уровни логов (debug, info, warn, error).

---

## Структура проекта

- `cmd/`: Основные исполняемые файлы для запуска сервиса и продюсера Kafka.
- `internal/`: Внутренние компоненты:
    - `cache/`: Кэширование данных в памяти.
    - `config/`: Парсинг конфигов.
    - `kafka/`: Логика обработки сообщений Kafka.
    - `model/`: Парсинг JSON файла.
    - `handlers/`: HTTP-сервер и обработка запросов.
    - `storage/`: Работа с базой данных PostgreSQL.
    - `logger/`: Логирование с помощью `slog`.
- `config/`: Config файл.
- `data/`: Модели json и sql.
- `scripts/`: Скрипт для ожидания подключения к Kafka и Postgres.
- `static/`: UI.
- `Dockerfile`: Файл для сборки Docker образа.
- `docker-compose.yml`: Настройка сервисов через Docker Compose.
- `Makefile`: Упрощение команд для сборки и запуска сервиса.
- `tests/`: Тесты.

---