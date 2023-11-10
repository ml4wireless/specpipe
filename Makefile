# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install

build: build-server build-edge
build-server:
	$(GOBUILD) -ldflags="-X main.Version=v0.3.2 -w -s" -o sp-server ./cmd/server
build-edge:
	$(GOBUILD) -ldflags="-X main.Version=v0.3.2 -w -s" -o sp-edge ./cmd/edge

codegen: codegen-install
	swagger-codegen generate -l openapi-yaml -i server/openapi/main.yaml -t server/openapi -DoutputFile=merge.yaml
	$(shell $(GOCMD) env GOPATH)/bin/oapi-codegen -package server -generate "types,gin,spec" merge.yaml > server/server.gen.go
	rm -rf merge.yaml .swagger-codegen/
codegen-install:
	$(GOINSTALL) github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.0.0

docker: docker-server docker-edge
docker-server:
	@docker build -f ./build/server/Dockerfile --build-arg VERSION=v0.3.2 -t minghsu0107/specpipe-server .
docker-edge:
	@docker build -f ./build/edge/Dockerfile --build-arg VERSION=v0.3.2 -t minghsu0107/specpipe-edge .
clean:
	$(GOCLEAN)
	rm -f sp-server sp-edge