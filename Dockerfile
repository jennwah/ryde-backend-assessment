FROM golang:1.23.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /ryde-backend cmd/service.go

FROM gcr.io/distroless/base

COPY --from=builder /ryde-backend /ryde-backend

EXPOSE 8080

CMD ["/ryde-backend"]
