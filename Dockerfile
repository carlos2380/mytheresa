FROM golang:1.22.4-alpine as builder

WORKDIR /app 
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /mytheresa cmd/server/main.go
RUN go build -o client cmd/client/main.go

FROM alpine:latest as mytheresa
WORKDIR /mytheresa
COPY --from=builder /app/docs /mytheresa/docs
COPY --from=builder /mytheresa .
RUN chmod -R 755 /mytheresa/docs
RUN chmod +x /mytheresa/mytheresa
ENV PORT=8000
EXPOSE 8000
CMD ["./mytheresa"]

FROM alpine:latest as client
WORKDIR /
COPY --from=builder /app/client /client
CMD ["/client"]
