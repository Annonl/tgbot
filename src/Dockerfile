FROM golang:alpine
WORKDIR /app
COPY . .
RUN go build -o bot main.go
CMD ["./bot"]