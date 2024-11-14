FROM golang:1.21-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o /app/server .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/server /server

EXPOSE 8080

CMD ["/server"]
