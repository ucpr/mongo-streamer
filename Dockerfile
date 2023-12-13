## Buidler
FROM golang:1.21.4-alpine3.18 AS builder

ENV VERSION=""
ENV REVISION=""
ENV TIMESTAMP=""

WORKDIR /app

COPY ./ ./

RUN go mod download && go build -ldflags "-X main.BuildVersion=${VERSION} -X main.BuildRevision=${REVISION} -X main.BuildTimestamp=${TIMESTAMP}" -o ./build/mongo-streamer ./cmd

## Runner
FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/build/mongo-streamer .

USER 1001

CMD [ "/app/mongo-streamer" ]
