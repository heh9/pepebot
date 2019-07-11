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

# Install project dependencies
RUN go get

# Expose port
EXPOSE 9001

# build project
CMD ["go", "run", "*.go"]