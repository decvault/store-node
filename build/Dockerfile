FROM golang:1.25.5-alpine AS build

WORKDIR /src

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal ./internal/

RUN go build -o ./bin/main ./cmd/

FROM alpine:latest AS app

WORKDIR /app

EXPOSE 50051

COPY config/ ./config/
COPY --from=build /src/bin/main .

CMD ["./main"]
