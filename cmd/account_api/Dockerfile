FROM golang:1.25.1

WORKDIR /app

COPY ../../go.mod ../../go.sum ./
RUN go mod download

COPY ../../. ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /account-api ./cmd/account_api

EXPOSE 9090

CMD ["/account-api"]
