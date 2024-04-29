FROM golang:1.22-alpine as builder
WORKDIR /app
RUN apk add --no-cache make nodejs npm
COPY . ./
COPY .env.production .env
RUN make docker-build-install

EXPOSE 4000
ENTRYPOINT ["./bin/web"]