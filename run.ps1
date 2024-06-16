# Run npm build in the resources directory
npm --prefix ./resources run build

# Run the Go application
go run .\main.go .\scan.go .\throttle.go .\ui.go .\hotkeys.go

