FROM golang:buster

WORKDIR /go/src/crawler
COPY . /go/src/crawler
RUN go mod download
RUN make build

CMD [ "/go/src/crawler/bin/p2pcrawler" ]