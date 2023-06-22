from golang:1.20-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:3.13

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]