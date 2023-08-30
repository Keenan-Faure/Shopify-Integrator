FROM debian:stable-slim

FROM golang:1.20
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN mkdir keenan && cd keenan
WORKDIR /keenan/

# Migrations
ADD /sql/ /keenan/sql/
ADD /scripts/migrations.sh /keenan/sql/schema
COPY .env /keenan/sql/schema/

RUN ["chmod", "+x", "/keenan/sql/schema/migrations.sh"]

ADD integrator /keenan/
CMD [ "/keenan/integrator" ]