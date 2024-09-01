$version = Get-Date -Format "yyyy-MM-dd HH:mm"
$BUILD_NAME = "kocli"

function Build {
    go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$version'" -o "./$BUILD_NAME" "./main.go"
}

function Mac {
    $env:GOOS = "darwin"
    go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$version'" -o "./$BUILD_NAME-darwin" "./main.go"
    if (Get-Command upx -ErrorAction SilentlyContinue) {
        upx "$BUILD_NAME-darwin"
    }
}

function Win {
    $env:GOOS = "windows"
    go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$version'" -o "./$BUILD_NAME.exe" "./main.go"
    if (Get-Command upx -ErrorAction SilentlyContinue) {
        upx "$BUILD_NAME.exe"
    }
}

function Linux {
    $env:GOOS = "linux"
    go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$version'" -o "./$BUILD_NAME-linux" "./main.go"
    if (Get-Command upx -ErrorAction SilentlyContinue) {
        upx "$BUILD_NAME-linux"
    }
}

function Build-Test {
    go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$version'" -o "./${BUILD_NAME}_debug_bin.exe" "./main.go"
}

function Init-Test {
    Build-Test
    Init-Example
}

function Init-Example {
    & "./${BUILD_NAME}_debug_bin" init -p "github.com/woocoos/helloworld" -t "project/internal/integration/helloworldtest"
    Set-Location "project/internal/integration/helloworldtest/cmd"
    go run "main.go"
}