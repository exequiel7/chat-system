FROM golang:1.22-alpine

ENV GO111MODULE=on

RUN apk add --no-cache git gcc g++ curl make
RUN wget https://github.com/jwilder/dockerize/releases/download/v0.6.1/dockerize-linux-amd64-v0.6.1.tar.gz
RUN tar -C /usr/local/bin -xzvf dockerize-linux-amd64-v0.6.1.tar.gz

WORKDIR /go/src/chat-system
COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080

CMD ["dockerize", "-wait", "tcp://cassandra:9042", "-timeout", "120s", "-wait-retry-interval", "15s", "go", "run", "cmd/api/main.go"]

