build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/greeter
deploy: build
	serverless deploy
remove:
	serverless remove
