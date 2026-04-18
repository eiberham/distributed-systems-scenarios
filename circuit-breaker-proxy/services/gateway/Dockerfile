FROM golang:1.23-alpine

ENV GOTOOLCHAIN=auto

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o gateway ./services/gateway

EXPOSE 8081

CMD ["./gateway"]