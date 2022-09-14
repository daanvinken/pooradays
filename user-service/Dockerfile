FROM golang:latest as builder

LABEL maintainer="Daan Vinken"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o user-service .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/
RUN mkdir config

COPY --from=builder /app/user-service .
COPY  config/config.dev.json /root/config

EXPOSE 4000

# Command to run the executable
CMD ["./user-service"]