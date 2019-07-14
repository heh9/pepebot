FROM golang:1.12

# Update and install curl
RUN apt-get update

# Creating work directory
RUN mkdir $GOPATH/src/pepe.bot

# Adding project to work directory
ADD . $GOPATH/src/pepe.bot

# Choosing work directory
WORKDIR $GOPATH/src/pepe.bot

RUN git clone https://go.googlesource.com/crypto $GOPATH/src/golang.org/x/crypto

# Install project dependencies
RUN go get -t

# Expose port
EXPOSE 9001

# build project
RUN go build -o pepe_bot .

CMD ["./pepe_bot"]