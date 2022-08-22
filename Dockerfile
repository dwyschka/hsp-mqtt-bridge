FROM golang:1.18-alpine AS builder
COPY ./app/* /src/
RUN mkdir /build
WORKDIR /src
RUN CGO_ENABLED=0 go build -o /build/hsp-stove main

FROM alpine:3.15
COPY --from=builder /build/hsp-stove /bin/hsp-stove
RUN chmod +x /bin/hsp-stove
ENTRYPOINT ["/bin/hsp-stove"]