NAME=ifwatch
VER=0.0.1
REL=1
SRC=$(shell glide name)

all: build

deps:
	go get -u golang.org/x/lint/golint
	go get -u github.com/kisielk/errcheck
	go get -u github.com/Masterminds/glide/...

proto:
	protoc -I network/ network/network.proto --go_out=plugins=grpc:network/

clean:
	rm -f $(NAME) $(NAME)-*.rpm $(NAME)-*.tar.gz $(NAME).1 $(NAME).toml.5
	which docker &>/dev/null && docker rm $(NAME) || true

format:
	gofmt -w .

test: format
	golint -set_exit_status .
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go vet .
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 errcheck .
#	go test . -v -covermode=atomic

man:
	pandoc man/${NAME}.1.md -s -t man -o ${NAME}.1
	pandoc man/${NAME}.toml.5.md -s -t man -o ${NAME}.toml.5

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(VER)"

rpm: clean
	tar -czf $(NAME)-$(VER)-$(REL).tar.gz .
	docker build --pull=true \
		--build-arg NAME=$(NAME) \
		--build-arg VER=$(VER) \
		--build-arg REL=$(REL) \
		--build-arg SRC=$(SRC) -t $(NAME):latest .
	cid=$$(docker create -ti --name $(NAME) $(NAME):latest true) ;\
	docker cp $$cid:/root/rpmbuild/RPMS/x86_64/$(NAME)-$(VER)-$(REL).x86_64.rpm . ;\
	docker cp $$cid:/root/rpmbuild/SRPMS/$(NAME)-$(VER)-$(REL).src.rpm .

.PHONY: all deps proto clean format test man build linux rpm
