## Buidler
FROM golang:1.21.4-alpine3.18 AS builder

WORKDIR /app

COPY ./ ./

RUN go mod download && go build -o ./build/mongo-streamer ./cmd/main.go

## Runner
FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/build/mongo-streamer .

USER 1001

CMD [ "/app/mongo-streamer" ]
