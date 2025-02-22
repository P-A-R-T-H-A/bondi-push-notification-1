FROM golang:1.22-alpine as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app

FROM gcr.io/distroless/base-debian11 as runner
COPY --from=builder /app /app
EXPOSE 8070
ENTRYPOINT ["/app/app"]


