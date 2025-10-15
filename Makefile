include make/lint.mk
include make/build.mk

lint: cart-lint loms-lint notifier-lint comments-lint

build: cart-build loms-build notifier-build comments-build

run-all:
	docker compose up --build

run-monitor:
	docker compose -f docker-compose-monitoring.yml up -d

stop-monitor:
	docker compose -f docker-compose-monitoring.yml down


run-kafka:
	docker compose -f docker-compose-kafka.yml up -d

stop-kafka:
	docker compose -f docker-compose-kafka.yml down


test:
	$(MAKE) -C ./cart test
	$(MAKE) -C ./loms test

test-cover:
	$(MAKE) -C ./cart test-cover
	$(MAKE) -C ./loms test-cover

test-race:
	$(MAKE) -C ./cart test-race
	$(MAKE) -C ./loms test-race


test-integration:
	$(MAKE) -C ./loms test-integration


test-api:
	rm -rf tests/tests/allure-results
	go test ./tests/... -tags=api

test-api-win:
	if exist .\tests\tests\allure-results rmdir /s /q .\tests\tests\allure-results
	go test ./tests/... -tags=api

allure:
	allure serve ./tests/tests/allure-results
