FROM golang:1.22-alpine AS builder

WORKDIR /go/src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gophermart ./cmd/gophermart


FROM scratch

COPY --from=builder /go/src/gophermart /
COPY ./migrations /migrations
CMD ["/gophermart"]
