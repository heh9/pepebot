FROM alpine AS builder

RUN apk update

RUN apk upgrade

RUN apk add --update go=1.16 gcc=6.3.0-r4 g++=6.3.0-r4

WORKDIR /build

ADD . /build

RUN CGO_ENABLED=1 go build -o pepebot .

FROM alpine

COPY --from=builder /build/pepebot /usr/bin/pepebot

# Expose port
EXPOSE 9001

ENTRYPOINT ["/usr/bin/pepebot"]
