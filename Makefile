BINARY_NAME=pangea

build:
	mkdir -p bin/
	GOARCH=amd64 GOOS=darwin go build -o bin/${BINARY_NAME}-darwin-x86_64
	GOARCH=amd64 GOOS=linux go build -o bin/${BINARY_NAME}-linux-x86_64
	GOARCH=amd64 GOOS=windows go build -o bin/${BINARY_NAME}-windows-x86_64.exe
	GOARCH=arm64 GOOS=darwin go build -o bin/${BINARY_NAME}-darwin-arm
	GOARCH=arm64 GOOS=linux go build -o bin/${BINARY_NAME}-linux-arm
	GOARCH=arm64 GOOS=windows go build -o bin/${BINARY_NAME}-windows-arm.exe

run: build
	./${BINARY_NAME}

dev:
	go run main.go

clean:
	rm -rf bin/*