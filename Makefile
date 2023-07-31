BINARY_NAME=pangea

build:
	mkdir -p bin/
	GOARCH=amd64 GOOS=darwin go build -o bin/${BINARY_NAME}-darwin-amd64
	GOARCH=amd64 GOOS=linux go build -o bin/${BINARY_NAME}-linux-amd64
	GOARCH=amd64 GOOS=windows go build -o bin/${BINARY_NAME}-windows-amd64.exe
	GOARCH=arm64 GOOS=darwin go build -o bin/${BINARY_NAME}-darwin-arm
	GOARCH=arm64 GOOS=linux go build -o bin/${BINARY_NAME}-linux-arm
	GOARCH=arm64 GOOS=windows go build -o bin/${BINARY_NAME}-windows-arm.exe

run: build
	./${BINARY_NAME}

dev:
	go run main.go

clean:
	rm -rf bin/*