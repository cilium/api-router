all:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' .
	docker build -t cilium/api-router .
