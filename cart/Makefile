BINDIR=${CURDIR}/bin
PACKAGE=route256/cart

bindir:
	mkdir -p ${BINDIR}

build: bindir
	echo "build cart"

minimock-gen:
	mkdir -p mocks
	$(MAKE) minimock-gen-common

minimock-gen-win:
	mkdir mocks 2>nul || exit 0
	$(MAKE) minimock-gen-common

minimock-gen-common:
	minimock -i route256/cart/internal/service.CartRepository -o ./mocks/ -s "_mock.go"
	minimock -i route256/cart/internal/service.ProductService -o ./mocks/ -s "_mock.go"
	minimock -i route256/cart/internal/handler.CartService -o ./mocks/ -s "_mock.go"

test:
	go test ./... -count=1

test-cover:
	go test ./... -cover

test-api:
	rm -rf ./tests/allure-results
	go test ./tests/... -tags=api

test-api-win:
	if exist tests\allure-results rmdir /s /q tests\allure-results
	go test ./tests/... -tags=api

allure: 
	allure serve ./tests/allure-results
