ARG  GOLANG_VERSION=1.11.4
FROM docker.io/library/golang:${GOLANG_VERSION} AS libzfs-depenencies
RUN  apt-get update
RUN  apt-get install -y zlib1g-dev uuid-dev libblkid-dev

FROM libzfs-depenencies AS libzfs-compiler
ARG  ZFS_VERSION=0.7.12
RUN  curl -sL https://github.com/zfsonlinux/zfs/releases/download/zfs-${ZFS_VERSION}/zfs-${ZFS_VERSION}.tar.gz \
     | tar zxC /usr/src

RUN  cd /usr/src/zfs-${ZFS_VERSION} \
 &&  ./configure --prefix=/usr --enable-static -with-config=user \
 &&  make -l8 install DESTDIR=/libzfs

FROM libzfs-depenencies AS compiler
RUN  apt-get install -y protobuf-compiler
RUN  go get github.com/golang/protobuf/protoc-gen-go
RUN  git clone https://github.com/grpc-ecosystem/grpc-gateway /go/src/github.com/grpc-ecosystem/grpc-gateway
COPY --from=libzfs-compiler /libzfs /
WORKDIR /go/src/github.com/steigr/zfs
COPY go.mod go.mod
COPY go.sum go.sum
ENV  GO111MODULE=on
RUN  go mod download
COPY Makefile Makefile
COPY cmd cmd
COPY pkg pkg
RUN  make proto
RUN  make zfsd

FROM docker.io/library/alpine:3.8 AS zfsd
COPY --from=compiler /go/bin/zfsd /bin/zfsd
ENTRYPOINT ["zfsd"]