FROM golang:1.23-alpine AS builder

WORKDIR /build

COPY ./notifier/go.mod ./notifier/go.mod
COPY ./notifier/go.sum ./notifier/go.sum
COPY ./cart ./cart
COPY ./loms ./loms

WORKDIR /build/notifier

RUN go mod tidy
RUN go mod download

COPY ./notifier .

RUN CGO_ENABLED=0 GOOS=linux go build -o /notifier ./cmd/main.go

FROM scratch
COPY --from=builder notifier /bin/notifier
COPY notifier/configs/values_local.yaml /bin/config/values_local.yaml

ENV CONFIG_FILE=/bin/config/values_local.yaml

ENTRYPOINT ["/bin/notifier"]
