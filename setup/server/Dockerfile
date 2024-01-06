FROM debian:stable-slim

FROM golang:1.20

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN mkdir keenan && cd keenan
WORKDIR /keenan/

# Migrations
COPY ./ /keenan/
COPY ./.env /keenan/

RUN ["chmod", "+x", "/keenan/scripts/migrations.sh"]

RUN cd /keenan/

CMD [ "/bin/sh", "-c" , "go build -o integrator && ./integrator"]