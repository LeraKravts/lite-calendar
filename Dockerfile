FROM golang:1.24-bullseye

RUN apt-get update && apt-get install -y git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd

CMD ["./app"]
