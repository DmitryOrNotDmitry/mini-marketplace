# mini-marketplace

-	Сервисы написаны на go, взаимодействуют по gRPC и REST.

-	PostgreSQL (драйвер – pgx/v5), репликация на master + sync, шардирование с согласованным хэшированием, миграции (goose), SQL-запросы реализованы с помощью sqlc (кодогенерации).

-	Unit-тесты и интеграционные (minimock, testify) тесты, e2e тесты (allure-go).

-	Логирование (uber-go/zap), метрики (Prometheus, Grafana), трассировка (Jaeger, OpenTelemetry).

-	Kafka (доступ через sarama), создана сематика at least once, для транзакционной согласованности (записи в PostgreSQL и доставки сообщения в Kafka) со стороны producer’а реализован паттерн outbox.

-	Docker, Docker Compose для локального запуска.

-	Для всех сервисов реализован graceful shutdown.

-	Реализован клиентский ratelimiter к внешнему сервису.

-	Сбор профилей через pprof в Makefile.
