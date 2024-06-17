# Run npm build in the resources directory
npm --prefix ./resources run build

# Run the Go application
$env:CGO_CFLAGS="-IC:/SDL2/include"
$env:CGO_LDFLAGS="-LC:/SDL2/lib/x64"
go run .\main.go .\scan.go .\throttle.go .\ui.go .\hotkeys.go .\sdl.go .\throttlecontroller.go .\joy1.go

