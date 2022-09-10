FROM golang:latest

RUN mkdir -p /usr/local/go/src/product-listing

WORKDIR /usr/local/go/src/product-listing

COPY . /usr/local/go/src/product-listing

RUN go install product-listing

CMD /usr/local/go/bin/product-listing

EXPOSE 8080
