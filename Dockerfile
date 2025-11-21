FROM golang:1.25.4-alpine3.22

WORKDIR /project/pr-service

COPY . .
RUN go mod tidy

WORKDIR /project/pr-service/cmd/main
RUN go build main.go && mkdir -p /project/data/logs

EXPOSE 8080
VOLUME /project/data

ENTRYPOINT ["./main"]
