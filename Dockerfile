FROM golang:alpine AS Builder
LABEL maintainer="Abdulsamet İleri <abdulsamet.ileri@ceng.deu.edu.tr>"
WORKDIR /app
COPY go.sum go.mod ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build --tags prod -o main main.go
FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env .
EXPOSE 8080
CMD ["./main"]
