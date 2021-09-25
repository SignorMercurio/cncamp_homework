export tag=1.0.0

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-w -s" -gcflags "-N -l" -o bin/httpserver

release:
	docker buildx build --platform linux/amd64 --push -t signormercurio/httpserver:${tag} .