## Buidler
FROM golang:1.21.4-alpine3.18 AS builder

ARG VERSION
ARG REVISION
ARG TIMESTAMP

WORKDIR /app

COPY ./ ./

RUN go mod download && go build -ldflags "-X github.com/ucpr/mongo-streamer/pkg/stamp.BuildVersion=${VERSION} -X github.com/ucpr/mongo-streamer/pkg/stamp.BuildRevision=${REVISION} -X github.com/ucpr/mongo-streamer/pkg/stamp.BuildTimestamp=${TIMESTAMP}" -o ./build/mongo-streamer ./cmd

## Runner
FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/build/mongo-streamer .

USER 1001

CMD [ "/app/mongo-streamer" ]
