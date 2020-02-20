FROM golang:1.13-alpine

ARG ARG_BOT_WEBHOOK
ARG ARG_BOT_TOKEN

ENV BOT_WEBHOOK=$ARG_BOT_WEBHOOK
ENV BOT_TOKEN=$ARG_BOT_TOKEN

RUN apk add --update git gcc libc-dev openssl

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

RUN openssl genrsa 2048 | openssl pkcs8 -topk8 -nocrypt -out private.pem

CMD ["app"]