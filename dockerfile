FROM golang:latest as builder

WORKDIR /api
COPY go.* ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o api

FROM alpine:3
RUN apk add --no-cache ca-certificates
COPY --from=builder /api/api /api

CMD ["/api"]
