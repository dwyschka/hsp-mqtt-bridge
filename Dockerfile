FROM alpine:3.15
WORKDIR /

COPY ./app/Stove /bin/Stove
RUN chmod +x /bin/Stove


ENTRYPOINT ["/bin/Stove"]