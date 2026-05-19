FROM golang:1.24 AS builder
WORKDIR /src
COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/evidence-api ./cmd/api

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /out/evidence-api /app/evidence-api
EXPOSE 8080
ENV APP_PORT=8080
ENV APP_API_KEY=dev-api-key
ENTRYPOINT ["/app/evidence-api"]
