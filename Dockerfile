FROM golang:1.16-alpine AS builder

# Creating work directory
RUN mkdir /build

# Adding project to work directory
ADD . /build

# Choosing work directory
WORKDIR /build

# build project
RUN go build -o pepebot .

FROM alpine

COPY --from=builder /build/pepebot /usr/bin/pepebot

# Expose port
EXPOSE 9001

ENTRYPOINT ["/usr/bin/pepebot"]
