package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/subhasundardass/tuix/internal/app"
	"github.com/subhasundardass/tuix/tuix"
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			// Write panic to a file
			f, _ := os.Create("/tmp/panic.log")
			fmt.Fprintf(f, "panic: %v\n\n", r)
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			fmt.Fprintf(f, "%s", buf[:n])
			f.Close()
		}
	}()

	tuiApp := tuix.NewApp(0, 0) // Full Screen

	tuiApp.Run(app.Root, tuix.Props{})
}
