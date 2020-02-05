FROM golang:1.13-alpine

ARG ARG_BOT_WEBHOOK
ARG ARG_BOT_TOKEN

ENV BOT_WEBHOOK=$ARG_BOT_WEBHOOK
ENV BOT_TOKEN=$ARG_BOT_TOKEN

RUN apk add --update git gcc libc-dev

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]