FROM golang:1.23.2

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /ryde-backend cmd/service.go

CMD ["/ryde-backend"]