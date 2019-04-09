FROM golang:1.11.5-stretch

WORKDIR $GOPATH/src/github.com/zjshen14/iotex-leader-election/

RUN apt-get install -y --no-install-recommends make

COPY . .

RUN make clean build && \
	ln -s $GOPATH/src/github.com/zjshen14/iotex-leader-election/bin/elector /usr/local/bin/elector

CMD ["elector"]