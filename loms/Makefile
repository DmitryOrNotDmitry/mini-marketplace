BINDIR=${CURDIR}/bin
LOCAL_BIN:=$(CURDIR)/bin
PACKAGE=route256/loms

PROTO_FILES_PATH := "api"

bindir:
	mkdir -p ${BINDIR}

build: bindir
	echo "build loms"

test:
	go test ./internal/... -count=1

test-cover:
	go test ./internal/... -count=1 -cover

test-race:
	go test ./internal/... -count=1 -race

test-integration:
	go test ./test/integration/... -tags=integration
	
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
	minimock -i route256/loms/internal/service.StockRepoFactory -o ./mocks/ -s "_mock.go"
	minimock -i route256/loms/internal/service.OrderRepoFactory -o ./mocks/ -s "_mock.go"
	minimock -i route256/loms/internal/handler.StockService -o ./mocks/ -s "_mock.go"
	minimock -i route256/loms/internal/handler.OrderService -o ./mocks/ -s "_mock.go"


bin-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 && \
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0 && \
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.0.4 && \
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.19.1 && \
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.19.1 && \
	GOBIN=$(LOCAL_BIN) go install github.com/go-swagger/go-swagger/cmd/swagger@v0.30.5


bin-deps-win:
	set "GOBIN=$(LOCAL_BIN)" && \
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 && \
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0 && \
	go install github.com/envoyproxy/protoc-gen-validate@v1.0.4 && \
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.19.1 && \
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.19.1

vendor-rm:
	rm -rf vendor-proto

vendor-rm-win:
	if exist vendor-proto rmdir /s /q vendor-proto


vendor-proto/google/protobuf-common:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor-proto/protobuf &&\
	cd vendor-proto/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout

vendor-proto/google/protobuf: vendor-proto/google/protobuf-common
	mkdir -p vendor-proto/google
	mv vendor-proto/protobuf/src/google/protobuf vendor-proto/google
	rm -rf vendor-proto/protobuf

vendor-proto/google/protobuf-win: vendor-proto/google/protobuf-common
	if not exist vendor-proto\google mkdir vendor-proto\google
	if exist vendor-proto\protobuf\src\google\protobuf move vendor-proto\protobuf\src\google\protobuf vendor-proto\google
	if exist vendor-proto\protobuf rmdir /s /q vendor-proto\protobuf


vendor-proto/validate-common:
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
	https://github.com/bufbuild/protoc-gen-validate vendor-proto/tmp && \
	cd vendor-proto/tmp && \
	git sparse-checkout set --no-cone validate &&\
	git checkout

vendor-proto/validate: vendor-proto/validate-common
	mkdir -p vendor-proto/validate
	mv vendor-proto/tmp/validate vendor-proto
	rm -rf vendor-proto/tmp

vendor-proto/validate-win: vendor-proto/validate-common
	if exist vendor-proto\tmp\validate move vendor-proto\tmp\validate vendor-proto
	if exist vendor-proto\tmp rmdir /s /q vendor-proto\tmp


vendor-proto/google/api-common:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor-proto/googleapis && \
 	cd vendor-proto/googleapis && \
	git sparse-checkout set --no-cone google/api && \
	git checkout

vendor-proto/google/api: vendor-proto/google/api-common
	mkdir -p  vendor-proto/google
	mv vendor-proto/googleapis/google/api vendor-proto/google
	rm -rf vendor-proto/googleapis

vendor-proto/google/api-win: vendor-proto/google/api-common
	if exist vendor-proto\googleapis\google\api move vendor-proto\googleapis\google\api vendor-proto\google
	if exist vendor-proto\googleapis rmdir /s /q vendor-proto\googleapis


vendor-proto/protoc-gen-openapiv2/options-common:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway vendor-proto/grpc-ecosystem && \
 	cd vendor-proto/grpc-ecosystem && \
	git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
	git checkout

vendor-proto/protoc-gen-openapiv2/options: vendor-proto/protoc-gen-openapiv2/options-common
	mkdir -p vendor-proto/protoc-gen-openapiv2
	mv vendor-proto/grpc-ecosystem/protoc-gen-openapiv2/options vendor-proto/protoc-gen-openapiv2
	rm -rf vendor-proto/grpc-ecosystem

vendor-proto/protoc-gen-openapiv2/options-win: vendor-proto/protoc-gen-openapiv2/options-common
	if not exist vendor-proto\protoc-gen-openapiv2 mkdir vendor-proto\protoc-gen-openapiv2
	move vendor-proto\grpc-ecosystem\protoc-gen-openapiv2\options vendor-proto\protoc-gen-openapiv2\options
	if exist vendor-proto\grpc-ecosystem rmdir /s /q vendor-proto\grpc-ecosystem


vendor-proto: vendor-rm vendor-proto/google/protobuf vendor-proto/validate vendor-proto/google/api vendor-proto/protoc-gen-openapiv2/options
vendor-proto-win: vendor-rm-win vendor-proto/google/protobuf-win vendor-proto/validate-win vendor-proto/google/api-win vendor-proto/protoc-gen-openapiv2/options-win


protoc-generate:
	mkdir -p pkg/api/orders/v1
	mkdir -p pkg/api/stocks/v1
	mkdir -p api/openapiv2
	protoc \
	-I ${PROTO_FILES_PATH} \
	-I vendor-proto \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
	--go_out pkg/api \
	--go_opt paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
	--go-grpc_out pkg/api \
	--go-grpc_opt paths=source_relative \
	--plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate \
    --validate_out="lang=go,paths=source_relative:pkg/api" \
	--plugin=protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway \
    --grpc-gateway_out pkg/api \
    --grpc-gateway_opt logtostderr=true  \
	--grpc-gateway_opt paths=source_relative \
	--grpc-gateway_opt generate_unbound_methods=true \
	--plugin=protoc-gen-openapiv2=$(LOCAL_BIN)/protoc-gen-openapiv2 \
    --openapiv2_out api/openapiv2 \
    --openapiv2_opt logtostderr=true \
	api/orders/v1/orders.proto \
	api/stocks/v1/stocks.proto && \
	go mod tidy

protoc-generate-win:
	if not exist pkg\api\orders\v1 mkdir pkg\api\orders\v1
	if not exist pkg\api\stocks\v1 mkdir pkg\api\stocks\v1
	if not exist api\openapiv2 mkdir api\openapiv2
	protoc \
	-I ${PROTO_FILES_PATH} \
	-I vendor-proto \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go.exe \
	--go_out pkg/api \
	--go_opt paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc.exe \
	--go-grpc_out pkg/api \
	--go-grpc_opt paths=source_relative \
	--plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate.exe \
    --validate_out="lang=go,paths=source_relative:pkg/api" \
	--plugin=protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway.exe \
    --grpc-gateway_out pkg/api \
    --grpc-gateway_opt logtostderr=true  \
	--grpc-gateway_opt paths=source_relative \
	--grpc-gateway_opt generate_unbound_methods=true \
	--plugin=protoc-gen-openapiv2=$(LOCAL_BIN)/protoc-gen-openapiv2.exe \
    --openapiv2_out api/openapiv2 \
    --openapiv2_opt logtostderr=true \
	--openapiv2_opt merge_file_name=loms,allow_merge=true \
	api/orders/v1/orders.proto \
	api/stocks/v1/stocks.proto && \
	go mod tidy


CI_MIGRATIONS_FOLDER := ./loms/migrations/
CI_DB_DSN := postgresql://loms-user:loms-password@postgres-master:5432/loms_db?sslmode=disable

run-migrations:
	goose -dir $(CI_MIGRATIONS_FOLDER) postgres "$(CI_DB_DSN)" up


LOCAL_MIGRATIONS_FOLDER := ./migrations/
LOCAL_DB_DSN := postgresql://loms-user:loms-password@localhost:5432/loms_db?sslmode=disable

migrate:
	goose -dir $(LOCAL_MIGRATIONS_FOLDER) postgres "$(LOCAL_DB_DSN)" up

migrate-down:
	goose -dir $(LOCAL_MIGRATIONS_FOLDER) postgres "$(LOCAL_DB_DSN)" down
