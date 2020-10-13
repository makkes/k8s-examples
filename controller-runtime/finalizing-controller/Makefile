.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -o ctrl main.go

.PHONY: docker-build
docker-build: build
	docker build -t makkes/finalizing-controller:latest .

.PHONY: docker-push
docker-push: docker-build
	docker push makkes/finalizing-controller:latest