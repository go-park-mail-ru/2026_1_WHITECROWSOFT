FROM golang:latest AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /noterian ./cmd/main

FROM alpine:latest AS run

COPY --from=build /noterian /noterian

WORKDIR /app
EXPOSE 8000
CMD ["/noterian"]
