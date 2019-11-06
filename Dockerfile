FROM golang:1.13

WORKDIR /go/src

RUN cd /go/src \
  && git clone https://github.com/mgartner/pg_flame.git \
  && cd pg_flame \
  && go build

ENTRYPOINT [ "pg_flame/pg_flame" ]
