FROM golang:1.18.3-alpine3.16 as builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build main.go

FROM alpine:3.16

COPY --from=builder /app /app

EXPOSE 8000

CMD ["/app/main"]
