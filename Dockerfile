FROM golang:alpine AS dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
EXPOSE 8084
CMD ["go", "run", "main.go", "serve"]


FROM golang:alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app main.go


FROM alpine:latest AS prod

WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 8084
CMD ["./app", "serve"]
