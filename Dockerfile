FROM golang:1.19

WORKDIR /usr/app

COPY src/go.mod src/go.sum ./
RUN go mod download && go mod verify

COPY ./src .
RUN go build -v -o /usr/local/bin/app

CMD ["/usr/local/bin/app"]