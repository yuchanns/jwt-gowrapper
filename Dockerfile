FROM golang:alpine

WORKDIR /app

COPY . /app/

RUN go version && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go mod download

CMD ["go", "test"]

