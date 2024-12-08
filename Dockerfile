FROM golang:1.22.4-alpine as builder

WORKDIR /app 
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /mytheresa cmd/server/main.go

FROM alpine:latest as mytheresa
WORKDIR /mytheresa
COPY --from=builder /mytheresa .
RUN chmod +x /mytheresa/mytheresa
ENV PORT=8000
EXPOSE 8000
CMD ["./mytheresa"]

