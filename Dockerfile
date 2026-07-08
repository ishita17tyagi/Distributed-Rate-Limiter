FROM golang:1.22-alpine AS build

WORKDIR /src

COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal

RUN go build -o /bin/server ./cmd/server

FROM alpine:3.20

WORKDIR /app

COPY --from=build /bin/server /app/server

EXPOSE 8080

CMD ["/app/server"]
