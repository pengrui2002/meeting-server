FROM golang:1.21-alpine

WORKDIR /app

RUN apk add --no-cache git

ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on

COPY server/ .

RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

COPY web/ /web/

EXPOSE 8080

CMD ["/server"]