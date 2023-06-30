version := $(shell /bin/date "+%Y-%m-%d %H:%M")
BUILD_NAME=kocli_debug_bin

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
