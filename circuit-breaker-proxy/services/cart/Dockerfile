FROM golang:1.23-alpine

ENV GOTOOLCHAIN=auto

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o cart ./services/cart

EXPOSE 8083

CMD ["./cart"]