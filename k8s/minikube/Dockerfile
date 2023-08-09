# Dockerfile
FROM golang:1.16-buster AS builder
RUN mkdir /src
ADD . /src
WORKDIR /src

RUN go env -w GO111MODULE=auto
RUN CGO_ENABLED=0 go build -o main .

FROM alpine:3.5

RUN echo -e  "http://mirrors.aliyun.com/alpine/v3.5/main\nhttp://mirrors.aliyun.com/alpine/v3.5/community" \
     >  /etc/apk/repositories &&  \
    apk update &&  \
    apk add tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Shanghai/Asia" > /etc/timezone && \
    apk del tzdata

WORKDIR /
COPY --from=builder /src/main /main
EXPOSE 3000
ENTRYPOINT ["/main"]