FROM golang:1.23-alpine AS builder

WORKDIR /build

COPY ./comments/go.mod ./comments/go.mod
COPY ./comments/go.sum ./comments/go.sum
COPY ./cart ./cart
COPY ./loms ./loms

WORKDIR /build/comments

RUN go mod tidy
RUN go mod download

COPY ./comments .

RUN CGO_ENABLED=0 GOOS=linux go build -o /comments ./cmd/main.go

FROM scratch
COPY --from=builder comments /bin/comments
COPY comments/configs/values_local.yaml /bin/config/values_local.yaml

ENV CONFIG_FILE=/bin/config/values_local.yaml

ENTRYPOINT ["/bin/comments"]
