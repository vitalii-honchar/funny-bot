FROM golang:1.19

WORKDIR /usr/src/app

COPY build/cmd /usr/src/app/app
COPY migrations/ /usr/src/app/migrations/

CMD ["/usr/src/app/app"]