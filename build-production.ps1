npm --prefix ./resources run build
go build -ldflags "-s -w -H windowsgui" -o ./build/sbqol.exe .\main.go .\scan.go .\throttle.go .\ui.go .\hotkeys.go .\sdl.go