FROM golang:1.23-alpine AS builder
WORKDIR /usr/local/src

COPY ["go.mod", "go.sum", "/"]
RUN go mod download

COPY ./ ./
RUN go build -o ./main ./cmd/main.go


FROM alpine:3.21 AS runner

COPY --from=builder /usr/local/src/main /bin/main

CMD ["/bin/main"]