module route256/loms

go 1.23.0

toolchain go1.24.2

require (
	github.com/gojuno/minimock/v3 v3.4.7
	github.com/stretchr/testify v1.11.1
	google.golang.org/grpc v1.75.0
	google.golang.org/protobuf v1.36.6
	gopkg.in/yaml.v3 v3.0.1
	route256/cart v0.0.0-00010101000000-000000000000
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
)

replace route256/cart => ../cart
