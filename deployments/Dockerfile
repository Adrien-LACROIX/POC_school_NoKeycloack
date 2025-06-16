FROM golang:1.21

COPY . .

RUN go mod tidy

EXPOSE 8080

CMD ["go", "run", "cmd/main.go"]
