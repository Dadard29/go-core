FROM golang:1.13-alpine

ARG ARG_BOT_WEBHOOK
ARG ARG_BOT_TOKEN
ARG ARG_USERNAME_DB
ARG ARG_PASSWORD_DB
ARG ARG_JWT_SECRET

ENV BOT_WEBHOOK=$ARG_BOT_WEBHOOK
ENV BOT_TOKEN=$ARG_BOT_TOKEN
ENV USERNAME_DB=$ARG_USERNAME_DB
ENV PASSWORD_DB=$ARG_PASSWORD_DB
ENV JWT_SECRET=$ARG_JWT_SECRET

RUN apk add --update git gcc libc-dev openssl

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

RUN openssl genrsa 2048 | openssl pkcs8 -topk8 -nocrypt -out private.pem

CMD ["app"]