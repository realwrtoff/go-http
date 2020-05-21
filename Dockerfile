FROM centos:centos7

RUN yum install -y make \
    && yum install -y git gcc \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" >> /etc/timezone

# golang
ENV GOPROXY=https://goproxy.io
RUN curl -OL https://dl.google.com/go/go1.14.1.linux-amd64.tar.gz && \
    tar -xzvf go1.14.1.linux-amd64.tar.gz && mv go /usr/local
ENV PATH=$PATH:/usr/local/go/bin
ENV GOROOT=/usr/local/go

# 添加代码和编译
RUN git clone https://github.com/realwrtoff/go-http.git \
    && cd go-http && git pull && make output

EXPOSE 7060

WORKDIR /go-http/output/go-http
CMD [ "bin/server", "-c", "configs/server.json" ]
