FROM golang:1.14.1-buster

RUN cd / && git clone https://github.com/edenhill/librdkafka.git \
    && cd librdkafka \
    && ./configure --prefix /usr \
    && make \
    && make install

ENV PKG_CONFIG_PATH=/usr/lib/pkgconfig/
ENV LD_LIBRARY_PATH=/usr/lib

CMD cd /go/github.com/petergood/codecollabbackend && go test ./broker -v