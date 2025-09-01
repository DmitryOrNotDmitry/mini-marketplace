FROM golang:1.23-alpine AS builder

WORKDIR /build

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /cart ./cmd/main.go

FROM scratch
COPY --from=builder cart /bin/cart
COPY configs/values_local.yaml /bin/config/values_local.yaml

ENV CONFIG_FILE=/bin/config/values_local.yaml

ENTRYPOINT ["/bin/cart"]
