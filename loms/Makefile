BINDIR=${CURDIR}/bin
LOCAL_BIN:=$(CURDIR)/bin
PACKAGE=route256/loms

PROTO_FILES_PATH := "api"

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
	minimock -i route256/loms/internal/service.OrderRepository -o ./mocks/ -s "_mock.go"
	minimock -i route256/loms/internal/service.StockServiceI -o ./mocks/ -s "_mock.go"

bin-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 && \
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

bin-deps-win:
	set "GOBIN=$(LOCAL_BIN)" && \
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 && \
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

vendor-rm:
	rm -rf vendor-proto

vendor-rm-win:
	if exist vendor-proto rmdir /s /q vendor-proto

vendor-proto/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor-proto/protobuf &&\
	cd vendor-proto/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p vendor-proto/google
	mv vendor-proto/protobuf/src/google/protobuf vendor-proto/google
	rm -rf vendor-proto/protobuf

vendor-proto/google/protobuf-win:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor-proto/protobuf &&\
	cd vendor-proto/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	if not exist vendor-proto\google mkdir vendor-proto\google
	if exist vendor-proto\protobuf\src\google\protobuf move vendor-proto\protobuf\src\google\protobuf vendor-proto\google
	if exist vendor-proto\protobuf rmdir /s /q vendor-proto\protobuf

vendor-proto: vendor-rm vendor-proto/google/protobuf
vendor-proto-win: vendor-rm-win vendor-proto/google/protobuf-win

protoc-generate:
	mkdir -p pkg/api/orders/v1
	mkdir -p pkg/api/stocks/v1
	protoc \
	-I ${PROTO_FILES_PATH} \
	-I vendor-proto \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
	--go_out pkg/api \
	--go_opt paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
	--go-grpc_out pkg/api \
	--go-grpc_opt paths=source_relative \
	api/orders/v1/orders.proto \
	api/stocks/v1/stocks.proto && \
	go mod tidy

protoc-generate-win:
	if not exist pkg\api\orders\v1 mkdir pkg\api\orders\v1
	if not exist pkg\api\stocks\v1 mkdir pkg\api\stocks\v1
	protoc \
	-I ${PROTO_FILES_PATH} \
	-I vendor-proto \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go.exe \
	--go_out pkg/api \
	--go_opt paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc.exe \
	--go-grpc_out pkg/api \
	--go-grpc_opt paths=source_relative \
	api/orders/v1/orders.proto \
	api/stocks/v1/stocks.proto && \
	go mod tidy

