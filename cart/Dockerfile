FROM golang:1.23-alpine AS builder

WORKDIR /build

COPY ./cart/go.mod ./cart/go.mod
COPY ./cart/go.sum ./cart/go.sum
COPY ./loms ./loms

WORKDIR /build/cart

RUN go mod tidy
RUN go mod download

COPY ./cart .

RUN CGO_ENABLED=0 GOOS=linux go build -o /cart ./cmd/main.go

FROM scratch
COPY --from=builder cart /bin/cart
COPY cart/configs/values_local.yaml /bin/config/values_local.yaml

ENV CONFIG_FILE=/bin/config/values_local.yaml

ENTRYPOINT ["/bin/cart"]
