FROM golang:1.24-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN go build -o /bin/server ./cmd/server

FROM alpine:3.20

WORKDIR /app

COPY --from=build /bin/server /app/server

COPY client_limits.json /app/client_limits.json

EXPOSE 8080

CMD ["/app/server"]