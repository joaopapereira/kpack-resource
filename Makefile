ifdef VERSION
docker_registry = joaopapereira/kpack-resource:$(VERSION)
else
docker_registry = joaopapereira/kpack-resource
endif

docker: bin
	docker build -t $(docker_registry) .

publish: docker
	docker push $(docker_registry)

test:
	go test -v ./...

fmt:
	find . -name '*.go' | while read -r f; do \
		gofmt -w -s "$$f"; \
	done

bin: clean
	mkdir bin
	env GOOS=linux GOARCH=amd64 go build -o bin/in ./cmd/in/main.go
	env GOOS=linux GOARCH=amd64 go build -o bin/out ./cmd/out/main.go
	env GOOS=linux GOARCH=amd64 go build -o bin/check ./cmd/check/main.go

clean:
	rm -r ./bin/

.DEFAULT_GOAL := docker

.PHONY: go-mod docker-build docker-push docker test fmt
