FROM golang:alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./

RUN go build -o myapp .

EXPOSE 8080

CMD ["./myapp"]
