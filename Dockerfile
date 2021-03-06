FROM golang:1.13.1-alpine

WORKDIR /app

COPY . .

CMD ["./tasks"]