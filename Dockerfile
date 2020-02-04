FROM golang:1.13-alpine

ARG ARG_API_KEY

ENV API_KEY=$ARG_API_KEY

RUN apk add --update git gcc libc-dev

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]