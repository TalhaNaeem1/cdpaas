FROM golang:1.17.7-alpine3.15 AS building-stage
RUN apk add --update --no-cache alpine-sdk bash ca-certificates \
      libressl \
      tar \
      git openssh openssl yajl-dev zlib-dev cyrus-sasl-dev openssl-dev build-base coreutils
WORKDIR /go_service
COPY . ./
RUN GOOS=linux go build -a -tags musl .
FROM golang:1.17.7-alpine3.15
COPY --from=building-stage /go_service/pipelineService /
EXPOSE 8081
CMD ["/pipelineService"]

