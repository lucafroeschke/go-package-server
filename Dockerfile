FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o go-package-server .

FROM alpine:latest

WORKDIR /app/

COPY --from=builder /app/go-package-server .

EXPOSE 8080

CMD ["./go-package-server"]