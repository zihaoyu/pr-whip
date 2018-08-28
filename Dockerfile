FROM alpine:latest

COPY bin/whip /usr/bin/whip
ENTRYPOINT ["/usr/bin/whip"]
