FROM golang:1.16.0 AS Builder
WORKDIR /messaging-service
COPY go.sum go.mod ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build --tags prod -o main main.go

FROM alpine:3.7
WORKDIR /messaging-service
COPY --from=Builder /messaging-service/ .
CMD ["/messaging-service/main"]
