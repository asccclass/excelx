FROM alpine

WORKDIR /app

COPY ./app /app
RUN apk add ca-certificates

USER root
ENV PORT=80

ENTRYPOINT ["/app/app"]
