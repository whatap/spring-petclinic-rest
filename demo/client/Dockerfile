FROM golang:1.22.1 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY main.go .

RUN CGO_ENABLED=0 GOOS=linux go build -o spring-petclinic-rest-client
FROM alpine:3

WORKDIR /app

COPY --from=builder /app/spring-petclinic-rest-client .

RUN ls /app

ENTRYPOINT ["/app/spring-petclinic-rest-client"]

