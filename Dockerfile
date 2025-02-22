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

#FROM golang:1.22
#
#WORKDIR /app
#COPY . .
#RUN go mod download
#
#RUN go install github.com/beego/bee/v2@v2.1.0
#
#RUN bee pack
#
#RUN tar -xf app.tar.gz -c /app
#
#EXPOSE 8080
#
#CMD ["/app/app"]
