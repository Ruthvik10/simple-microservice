FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ ./

RUN go build -o ./event-bus-service

EXPOSE 3005

CMD ["./event-bus-service"]