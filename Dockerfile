FROM golang:1.14

# Update and install curl
RUN apt-get update

# Creating work directory
RUN mkdir /code

# Adding project to work directory
ADD . /code

# Choosing work directory
WORKDIR /code

# build project
RUN go build -o pepe_bot .

# Expose port
EXPOSE 9001

CMD ["./pepe_bot"]