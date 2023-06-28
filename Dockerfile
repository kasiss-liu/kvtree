FROM alpine:latest

LABEL author="kasiss"

COPY dist /app
ENV CNFPATH=/app/conf.toml

CMD ["/bin/sh","-c","cd /app/ && ./kvserver --config=${CNFPATH}"]
