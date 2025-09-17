FROM golang:1.24.4 AS builder

WORKDIR /app

COPY . .

RUN go build -o go-final .

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/go-final .
COPY --from=builder /app/web ./web
ENV TODO_PORT=7540
EXPOSE 7540

CMD ["./go-final"]