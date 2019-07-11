FROM golang:1.12

# Update and install curl
RUN apt-get update

# Creating work directory
RUN mkdir $GOPATH/src/pepe.bot

# Adding project to work directory
ADD . $GOPATH/src/pepe.bot

# Choosing work directory
WORKDIR $GOPATH/src/pepe.bot

# Generate .env file
ADD .env.example .env

RUN git clone https://go.googlesource.com/crypto $GOPATH/src/golang.org/x/crypto

# Install project dependencies
RUN go get -t

# Expose port
EXPOSE 9001

# build project
CMD ["go", "run", "*.go"]