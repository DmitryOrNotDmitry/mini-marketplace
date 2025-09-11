include make/lint.mk
include make/build.mk

lint: cart-lint loms-lint notifier-lint comments-lint

build: cart-build loms-build notifier-build comments-build

run-all:
	docker compose up --build

test:
	$(MAKE) -C ./cart test
	$(MAKE) -C ./loms test

test-cover:
	$(MAKE) -C ./cart test-cover
	$(MAKE) -C ./loms test-cover


test-api:
	rm -rf tests/tests/cart/allure-results
	ALLURE_OUTPUT_FOLDER=allure-results ALLURE_OUTPUT_PATH=/tests go test ./tests/...

test-api-win:
	if exist .\tests\tests\cart\allure-results rmdir /s /q .\tests\tests\cart\allure-results
	go test ./tests/...

allure:
	allure serve ./tests/tests/cart/allure-results
