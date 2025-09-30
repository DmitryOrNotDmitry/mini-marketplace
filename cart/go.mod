module route256/cart

go 1.23.0

toolchain go1.24.2

require (
	github.com/go-playground/assert/v2 v2.2.0
	github.com/gojuno/minimock/v3 v3.4.7
	github.com/stretchr/testify v1.11.1
	go.uber.org/goleak v1.3.0
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.75.1
	gopkg.in/yaml.v3 v3.0.1
	route256/loms v0.0.0-00010101000000-000000000000
)

require (
	github.com/envoyproxy/protoc-gen-validate v1.2.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.2 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/crypto v0.41.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250818200422-3122310a409c // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250818200422-3122310a409c // indirect
	google.golang.org/protobuf v1.36.7 // indirect
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-playground/validator/v10 v10.27.0
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
)

replace route256/loms => ../loms
