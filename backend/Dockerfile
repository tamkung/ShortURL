FROM golang:1.21.1-alpine3.17

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go mod tidy

EXPOSE 8000

CMD ["go", "run", "main.go"]
