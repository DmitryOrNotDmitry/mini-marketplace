BINDIR=${CURDIR}/bin
PACKAGE=route256/loms

bindir:
	mkdir -p ${BINDIR}

build: bindir
	echo "build loms"

run-migrations:
	echo "run migrations"

test:
	go test ./... -count=1

test-cover:
	go test ./... -count=1 -cover
	
minimock-gen:
	mkdir -p mocks
	$(MAKE) minimock-gen-common

minimock-gen-win:
	mkdir mocks 2>nul || exit 0
	$(MAKE) minimock-gen-common

minimock-gen-common:
	minimock -i route256/loms/internal/service.StockRepository -o ./mocks/ -s "_mock.go"
