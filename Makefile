version := $(shell /bin/date "+%Y-%m-%d %H:%M")
BUILD_NAME=kocli

build:
	go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" -o ./$(BUILD_NAME) ./main.go
mac:
	GOOS=darwin go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" -o ./$(BUILD_NAME)-darwin ./main.go
	$(if $(shell command -v upx), upx $(BUILD_NAME)-darwin)
win:
	GOOS=windows go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" -o ./$(BUILD_NAME).exe ./main.go
	$(if $(shell command -v upx), upx $(BUILD_NAME).exe)
linux:
	GOOS=linux go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" -o ./$(BUILD_NAME)-linux ./main.go
	$(if $(shell command -v upx), upx $(BUILD_NAME)-linux)

build-test:
	go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" -o ./$(BUILD_NAME)_debug_bin ./main.go

init-test: build-test init-example
init-example:
	./$(BUILD_NAME)_debug_bin init -p github.com/woocoos/helloworld -t project/internal/integration/helloworldtest
	cd project/internal/integration/helloworldtest/cmd && go run main.go