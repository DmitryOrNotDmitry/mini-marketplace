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
	minimock -i route256/cart/internal/service.HTTPClient -o ./mocks/ -s "_mock.go"
	minimock -i route256/cart/internal/service.LomsService -o ./mocks/ -s "_mock.go"
	minimock -i route256/cart/internal/service.RateLimiter -o ./mocks/ -s "_mock.go"
	minimock -i route256/cart/internal/handler.CartService -o ./mocks/ -s "_mock.go"
	minimock -i route256/cart/internal/handler.OrderCheckouter -o ./mocks/ -s "_mock.go"
	minimock -i route256/loms/pkg/api/stocks/v1.StockServiceV1Client -o ./mocks/ -s "_mock.go"
	minimock -i route256/loms/pkg/api/orders/v1.OrderServiceV1Client -o ./mocks/ -s "_mock.go"

test:
	go test ./internal/... -count=1

test-cover:
	go test ./internal/... -count=1 -cover

test-race:
	go test ./internal/... -count=1 -race

lint:
	gocyclo -over 10 -ignore "_test|mocks|tests|pkg|vendor" .
	gocognit -over 12 -ignore "_test|mocks|tests|pkg|vendor" .

lint-top:
	gocyclo -top 10 -ignore "_test|mocks|tests|pkg|vendor" .
	gocognit -top 10 -ignore "_test|mocks|tests|pkg|vendor" .
