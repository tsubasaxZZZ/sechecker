ARG GO_VERSION=1.14.5
#FROM golang:${GO_VERSION} AS build-stage
FROM golang:${GO_VERSION}-alpine AS build-stage
WORKDIR /src
COPY ./go.mod ./
RUN apk add make
RUN go mod download
COPY . .
RUN make build

FROM alpine:3.9
RUN apk add --no-cache dumb-init curl && mkdir /sechecker
COPY --from=build-stage /src/sechecker /sechecker/sechecker
COPY ./k8s-sample/crontab /var/spool/cron/crontabs/root
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD crond -l 1 -f