FROM ubuntu

RUN apt-get update &&\
    apt-get -y upgrade

RUN apt-get install -y curl nano

RUN curl https://getcaddy.com | bash -s personal http.cors,http.grpc

# Check caddy version
RUN caddy -version

# Create work directory
RUN mkdir /pepebot.caddy

# Select work directory
WORKDIR /pepebot.caddy

# Copy project to destination
ADD . /pepebot.caddy

# Expose ports
EXPOSE 80
EXPOSE 443

# Run project with caddy
ENTRYPOINT ["caddy"]
CMD ["-conf", "Caddyfile"]