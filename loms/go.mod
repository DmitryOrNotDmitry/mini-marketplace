module route256/loms

go 1.23.0

toolchain go1.24.2

require (
	github.com/envoyproxy/protoc-gen-validate v1.2.1
	github.com/gojuno/minimock/v3 v3.4.7
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.2
	github.com/jackc/pgx/v5 v5.7.6
	github.com/lib/pq v1.10.9
	github.com/stretchr/testify v1.11.1
	go.uber.org/goleak v1.3.0
	google.golang.org/genproto/googleapis/api v0.0.0-20250818200422-3122310a409c
	google.golang.org/grpc v1.75.1
	google.golang.org/protobuf v1.36.7
	gopkg.in/yaml.v3 v3.0.1
	route256/cart v0.0.0-00010101000000-000000000000
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.41.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250818200422-3122310a409c // indirect
)

replace route256/cart => ../cart
