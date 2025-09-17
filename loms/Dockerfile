FROM golang:1.23-alpine AS builder

WORKDIR /build

COPY ./loms/go.mod ./loms/go.mod
COPY ./loms/go.sum ./loms/go.sum
COPY ./cart ./cart

WORKDIR /build/loms

RUN go mod tidy
RUN go mod download

COPY ./loms .

RUN CGO_ENABLED=0 GOOS=linux go build -o /loms ./cmd/main.go

FROM scratch
COPY --from=builder loms /bin/loms
COPY loms/configs/values_local.yaml /bin/config/values_local.yaml
COPY loms/migrations /bin/migrations

ENV CONFIG_FILE=/bin/config/values_local.yaml

ENTRYPOINT ["/bin/loms"]
