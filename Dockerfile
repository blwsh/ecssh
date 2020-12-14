FROM golang:1.14-alpine AS builder
WORKDIR /src
COPY . /src
RUN go build .

FROM alpine:latest
RUN apk --no-cache add ca-certificates openssh-client && \
    addgroup -g 1000 app && \
    adduser -h /app -s /bin/sh -G app -S -u 1000 app
USER app
COPY --from=builder /src/ecssh /app/ecssh
ENTRYPOINT ["/app/ecssh"]
