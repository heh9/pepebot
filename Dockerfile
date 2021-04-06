FROM golang:1.16-alpine AS builder

RUN apk add build-base

# Choosing work directory
WORKDIR /build

# Adding project to work directory
ADD . /build

# build project
RUN go build -o pepebot .

FROM alpine

COPY --from=builder /build/pepebot /usr/bin/pepebot

# Expose port
EXPOSE 9001

ENTRYPOINT ["/usr/bin/pepebot"]
