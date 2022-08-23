FROM golang:1.18-alpine AS builder
RUN apk add upx
COPY ./app/* /src/
RUN mkdir /build
WORKDIR /src
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /build/hsp-stove main
RUN upx --brute /build/hsp-stove

FROM alpine:3.15
COPY --from=builder /build/hsp-stove /bin/hsp-stove
RUN chmod +x /bin/hsp-stove
ENTRYPOINT ["/bin/hsp-stove"]