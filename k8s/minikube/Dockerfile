FROM golang:1.16-buster AS builder
RUN mkdir /src
ADD . /src
WORKDIR /src

RUN go env -w GO111MODULE=auto
RUN go build -o main -ldflags "-w -extldflags -static"

FROM centos:6

#ENV TZ=America/Adak
ENV TZ=Asia/Shanghai

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN echo -e  "http://mirrors.aliyun.com/alpine/v3.5/main\nhttp://mirrors.aliyun.com/alpine/v3.5/community" \
     >  /etc/apk/repositories &&  \
    apk update &&  \
    apk add tzdata && \
    ln -sf /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone

WORKDIR /
COPY --from=builder /src/main /main
EXPOSE 3000
ENTRYPOINT ["main"]