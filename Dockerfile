FROM golang:1.22.5

WORKDIR /app

COPY . .

RUN go build -o ads-service cmd/main.go

EXPOSE 8080

CMD ["./ads-service"]