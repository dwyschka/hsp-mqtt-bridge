FROM golang:1.18-alpine AS builder
RUN apk add upx
COPY ./app/* /src/
RUN mkdir /build
WORKDIR /src
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /build/hsp-stove main
RUN upx --brute /build/hsp-stove

ARG BUILD_FROM
FROM $BUILD_FROM

COPY --from=builder /build/hsp-stove /bin/hsp-stove
RUN chmod +x /bin/hsp-stove

COPY run.sh /
RUN chmod a+x /run.sh

CMD [ "/run.sh" ]
