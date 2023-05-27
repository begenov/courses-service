FROM golang:alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o courses-services ./cmd/main.go

CMD ["./courses-services"]
