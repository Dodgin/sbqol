package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	"github.com/lxn/win"
)

func UiBootstrap(messageChannel chan string, doneChannel chan bool) {
	// Initialize astilectron
	ast, err := astilectron.New(nil, astilectron.Options{
		AppName: "SBQOL",
	})
	if err != nil {
		log.Fatalf("creating astilectron failed: %v", err)
	}
	defer ast.Close()

	// Start astilectron
	if err = ast.Start(); err != nil {
		log.Fatalf("starting astilectron failed: %v", err)
	}

	// Create a window
	window, err := ast.NewWindow(filepath.Join("resources", "index.html"), &astilectron.WindowOptions{
		Center:      astikit.BoolPtr(false),
		X:           astikit.IntPtr(int(win.GetSystemMetrics(win.SM_CXSCREEN)) - 256),
		Y:           astikit.IntPtr(32),
		Height:      astikit.IntPtr(256),
		Width:       astikit.IntPtr(256),
		Frame:       astikit.BoolPtr(false),
		AlwaysOnTop: astikit.BoolPtr(true),
		Resizable:   astikit.BoolPtr(false),
		Movable:     astikit.BoolPtr(true),
		Minimizable: astikit.BoolPtr(false),
		Maximizable: astikit.BoolPtr(false),
		Closable:    astikit.BoolPtr(true),
		SkipTaskbar: astikit.BoolPtr(true),
	})
	if err != nil {
		log.Fatalf("new window failed: %v", err)
	}

	// Show the window
	if err = window.Create(); err != nil {
		log.Fatalf("creating window failed: %v", err)
	}

	// Listen for messages from the main thread
	go func() {
		for {
			select {
			case msg := <-messageChannel:
				log.Printf("Received message from main thread: %s", msg)
				// Handle message (e.g., update UI)
				// ...

			case <-doneChannel:
				log.Println("Received done signal, exiting UI thread")
				return
			}
		}
	}()

	// Blocking pattern
	ast.Wait()

	// Signal the main thread that the UI is done
	doneChannel <- true
}

// Asset loads and returns the asset for the given name.
func Asset(name string) ([]byte, error) {
	path := filepath.Join("resources", name)
	return os.ReadFile(path)
}

// AssetDir returns the file names of the assets for the given name.
func AssetDir(name string) ([]string, error) {
	path := filepath.Join("resources", name)
	fileInfo, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, info := range fileInfo {
		names = append(names, info.Name())
	}
	return names, nil
}
