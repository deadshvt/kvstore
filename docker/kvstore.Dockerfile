FROM golang:alpine as builder

WORKDIR /project

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o bin/kvstore ./cmd/kvstore

FROM alpine:latest

WORKDIR /project

COPY .env .env
COPY --from=builder /project/bin/kvstore /project/bin/kvstore

CMD ["/project/bin/kvstore"]