FROM golang:1.23-alpine

ENV GOTOOLCHAIN=auto

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o orders ./services/orders

EXPOSE 8083

CMD ["./orders"]