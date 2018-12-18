all: proto zfsd

install: zfsd
	install -D -m 0755 $(GOPATH)/bin/zfsd bin/zfsd

proto: pkg/proto/zfs/zfs.pb.go pkg/proto/zpool/zpool.pb.go

pkg/proto/zfs/zfs.pb.go: pkg/proto/zfs.proto
	cd pkg/proto/zfs && protoc -I/usr/local/include -I../ \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:. \
		../zfs.proto

pkg/proto/zpool/zpool.pb.go: pkg/proto/zpool.proto
	cd pkg/proto/zpool && protoc -I/usr/local/include -I../ \
			-I${GOPATH}/src \
			-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:. \
		../zpool.proto

zfsd: proto zfsd.real
	@true

zfsd.real: /usr/lib/libzfs.a /usr/lib/libzpool.a /usr/lib/libnvpair.a /usr/lib/libzfs_core.a /usr/lib/libuutil.a /usr/lib/x86_64-linux-gnu/libz.a /usr/lib/x86_64-linux-gnu/libblkid.a /usr/lib/x86_64-linux-gnu/libpthread.a /usr/lib/x86_64-linux-gnu/libuuid.a /usr/lib/x86_64-linux-gnu/libm.a /usr/lib/x86_64-linux-gnu/libc.a
	GO111MODULE=on go get -ldflags "-linkmode external -extldflags '-static $?'" github.com/steigr/zfsd/cmd/zfsd

.PHONY: install
