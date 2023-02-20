FROM docker.io/golang:1.18.0-alpine3.14 AS builder

RUN apk fix
RUN apk --no-cache --update upgrade && \
    apk --no-cache --update add git git-lfs less openssh make && \
    git lfs install && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/*

ARG CI_COMMIT_ID
ARG CI_COMMIT_SHA
ARG BUILD_NUMBER

RUN mkdir -p /data/app/oar
COPY . /data/app/oar
WORKDIR /data/app/oar

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

RUN make all
RUN ls
WORKDIR /data/app/oar/bin
RUN ls

FROM docker.io/alpine:3.14.4
RUN mkdir -p /data/app/oar
WORKDIR /data/app/oar

COPY --from=builder /data/app/oar/bin .

ENTRYPOINT ["./oar"]
