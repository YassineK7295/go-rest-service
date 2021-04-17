FROM golang:alpine AS builder

RUN apk update ** apk add --no-cache github

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/membership-service/

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env .
# Skeleton required for viper to work with env vars
COPY --from=builder /app/config/skeleton.yaml ./local.yaml

EXPOSE 8080

CMD ["./main"]
