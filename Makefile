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
