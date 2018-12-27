DOCKERIMAGE ?= zfsd
STATIC_LIBS ?= /usr/lib/libzfs.a /usr/lib/libzpool.a /usr/lib/libnvpair.a /usr/lib/libzfs_core.a /usr/lib/libuutil.a /usr/lib/x86_64-linux-gnu/libz.a /usr/lib/x86_64-linux-gnu/libblkid.a /usr/lib/x86_64-linux-gnu/libpthread.a /usr/lib/x86_64-linux-gnu/libuuid.a /usr/lib/x86_64-linux-gnu/libm.a
PROTOCOLS    = $(wildcard pkg/proto/*.proto)
DEBUG_ADDRESS ?= 172.16.167.134

all: protocol zfsd

install: zfsd
	install -D -m 0755 $(GOPATH)/bin/zfsd bin/zfsd

protocol: $(PROTOCOLS)
	$(foreach var, \
		$?,protoc -I/usr/local/include -I$(dir $(var)) \
			   -I${GOPATH}/src \
			   -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
			   --go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:$(dir $(var)) \
			   $(var);)

zfsd: protocol zfsd.real
	@true

zfsd.real: $(STATIC_LIBS)
	go get -v -ldflags "-linkmode external -extldflags '-static $?'" github.com/steigr/zfsd/cmd/zfsd

.PHONY: install protocol

dockerimage:
	docker build -t $(DOCKERIMAGE) .

fmt:
	go fmt ./...

remote-debug: protocol dockerimage
	docker create --name=zfsd-copy-out $(DOCKERIMAGE)
	docker cp zfsd-copy-out:/bin/zfsd .
	docker rm zfsd-copy-out
	echo "pidof zfsd|xargs -r -t kill" | ssh root@$(DEBUG_ADDRESS) -- sh
	scp zfsd root@$(DEBUG_ADDRESS):/usr/local/sbin/zfsd
	ssh -tt root@$(DEBUG_ADDRESS) -- /usr/local/sbin/zfsd --logtostderr --v=9