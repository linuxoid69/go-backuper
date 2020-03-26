FROM alpine:3.9.5

COPY go-backuper /opt/

RUN apk update && \
    apk add postgresql-client

CMD ["/opt/go-backuper"]
