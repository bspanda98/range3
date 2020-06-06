# Build range3 in a stock Go builder container
FROM golang:1.13.7-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers

ADD . /range3
RUN cd /range3 && make geth

# Pull range3 into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /range3/build/bin/range3 /usr/local/bin/

EXPOSE 39796 39795 39797 39797/udp
ENTRYPOINT ["range3"]
