FROM golang:1.14-alpine AS build

COPY . /go/src/
WORKDIR /go/src/
RUN apk add make
RUN apk add git
RUN make build

FROM docker:20.10.1-dind

RUN apk update \
    && apk add tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

COPY --from=build /go/src/build /
WORKDIR /go-http
CMD [ "bin/server", "-c", "configs/server.json" ]
