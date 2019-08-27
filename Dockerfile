FROM alpine:3.10

RUN apk update \
    && apk add ca-certificates tzdata \
    && update-ca-certificates \
    && adduser -D -g '' eventhorizon

COPY . /opt/acesso

RUN cd /opt/acesso \
    && mkdir -p bin \
    && mv eventhorizon bin/

USER eventhorizon

CMD ["/opt/acesso/bin/eventhorizon"]
