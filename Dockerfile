FROM golang:1.20-alpine AS build

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o srv .

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/srv /app/

EXPOSE 8080

CMD ["/app/srv"]