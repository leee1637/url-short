FROM golang:1.26-alpine AS gobuild
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/app ./cmd/url-short



FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=gobuild /app/bin/app /app/app
COPY config /app/config/
EXPOSE 8080
CMD ["/app/app"]
