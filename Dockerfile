FROM golang:1.19

WORKDIR /usr/src/app

COPY build/cmd /usr/src/app/app

CMD ["/usr/src/app/app"]