ARG  ALPINE_VERSION=3.9
ARG  GOLANG_VERSION=1.11.4
FROM docker.io/library/alpine:${ALPINE_VERSION} AS alpine
FROM docker.io/library/golang:${GOLANG_VERSION} AS libzfs-dependencies
RUN  apt-get update
RUN  apt-get install -y zlib1g-dev uuid-dev libblkid-dev file dpkg-dev libssl-dev

FROM libzfs-dependencies AS libzfs-compiler
ARG  ZFS_VERSION=0.8.0-rc3
RUN  curl -sL https://github.com/zfsonlinux/zfs/releases/download/zfs-${ZFS_VERSION}/zfs-${ZFS_VERSION}.tar.gz \
     | tar zxC /usr/src

RUN  cd /usr/src/zfs-* \
 &&  ./configure --prefix=/usr --enable-static --with-config=user \
 &&  make -l8 install DESTDIR=/libzfs

FROM libzfs-dependencies AS compiler
RUN  apt-get install -y unzip
ARG  PROTOCOL_BUFFERS_VERSION=3.6.1
RUN  curl -L -o /tmp/protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOCOL_BUFFERS_VERSION}/protoc-${PROTOCOL_BUFFERS_VERSION}-linux-x86_64.zip \
 &&  cd /usr/local \
 &&  unzip /tmp/protoc.zip \
 &&  rm /tmp/protoc.zip
RUN  go get github.com/golang/protobuf/protoc-gen-go
RUN  git clone https://github.com/grpc-ecosystem/grpc-gateway /go/src/github.com/grpc-ecosystem/grpc-gateway
COPY --from=libzfs-compiler /libzfs /
WORKDIR /go/src/github.com/steigr/zfsd
COPY Makefile Makefile
ENV GO111MODULE=on
COPY go.mod go.mod
COPY go.sum go.sum
COPY cmd cmd
COPY pkg pkg
RUN  go mod vendor

# ZFS 0.8
RUN  rm -r vendor/github.com/bicomsystems/go-libzfs \
 &&  git clone https://github.com/steigr/go-libzfs --branch feature/zfs-0.8 vendor/github.com/bicomsystems/go-libzfs

RUN  GO111MODULE=off make zfsd

FROM alpine AS zfsd
COPY --from=compiler /go/bin/zfsd /bin/zfsd
ENTRYPOINT ["zfsd"]
