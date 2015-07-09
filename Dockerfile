FROM alpine:3.2

RUN apk --update add ca-certificates

COPY bunny /bunny

ENTRYPOINT ["/bunny"]