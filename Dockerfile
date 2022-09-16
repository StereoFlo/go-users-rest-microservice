# Compile stage
FROM golang:1.19 AS build-env

ADD . /var/www
WORKDIR /var/www

RUN go build -o /app

# Final stage
FROM debian:buster

EXPOSE 8000

WORKDIR /
COPY --from=build-env /app /

CMD ["/app"]