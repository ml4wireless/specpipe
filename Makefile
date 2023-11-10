# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean

all: build-server build-edge
build-server:
	$(GOBUILD) -ldflags="-X main.Version=v0.3.2 -w -s" -o sp-server ./cmd/server
build-edge:
	$(GOBUILD) -ldflags="-X main.Version=v0.3.2 -w -s" -o sp-edge ./cmd/edge

docker: docker-server docker-edge
docker-server:
	@docker build -f ./build/server/Dockerfile --build-arg VERSION=v0.3.2 -t minghsu0107/specpipe-server .
docker-edge:
	@docker build -f ./build/edge/Dockerfile --build-arg VERSION=v0.3.2 -t minghsu0107/specpipe-edge .
clean:
	$(GOCLEAN)
	rm -f sp-server sp-edge